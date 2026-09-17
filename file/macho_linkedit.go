package file

// The __LINKEDIT tables, generated from the image instead of patched into
// a template: an empty chained-fixups header, the exports trie, function
// starts, the symbol table, and the string table. The byte layouts follow
// dyld's formats (loader.h); the arrangement - regions in this order, each
// padded to 8, the trie node order of macho_trie.go - reproduces what ld
// emits for the reference two-symbol image byte for byte (the golden test
// pins that) and stays valid for any symbol set.

import (
	"encoding/binary"
	"slices"
)

// machoLinkedit - the generated tables and where they land in the file.
type machoLinkedit struct {
	data     []byte // fixups + trie + fstarts + nlists + strings
	trieAt   uint32 // offsets inside data
	fstartAt uint32
	symAt    uint32
	strAt    uint32
	strSize  uint32
	nsyms    uint32
	leOff    uint64 // absolute file offset (= place.total, page aligned)
	leSize   uint64 // len(data) + the signature length
	sigOff   uint64
	sigLen   uint64
	execSeg  uint32 // __TEXT vmsize: the CodeDirectory exec segment
}

// newMachOLinkedit generates the tables for a placed image with its
// resolved symbol addresses. sectByName maps the symbol section names to
// input section indices.
func newMachOLinkedit(
	place *machoPlacement,
	sectByName map[string]int,
	syms []MachOSym,
	addrs []uint64,
) machoLinkedit {
	var le machoLinkedit

	// section numbers for the nlists: 1-based, segments in file order
	sectNum := make([]uint32, len(place.sect))
	n := uint32(0)
	for si := range place.segs {
		for _, sp := range place.segs[si].sects {
			n++
			sectNum[sp.idx] = n
		}
	}

	// the exports trie: __mh_execute_header plus the globals, sorted
	names := []string{"__mh_execute_header"}
	globalAt := make(map[string]uint64, len(syms))
	for i := range syms {
		if syms[i].Global {
			globalAt[syms[i].Name] = addrs[i]
		}
	}

	for name := range globalAt {
		names = append(names, name)
	}
	slices.Sort(names)

	trie := machoNewTrie(names, func(name string) uint64 {
		if name == "__mh_execute_header" {
			return 0
		}

		return globalAt[name] - machoVMAddr
	})

	// function starts: every __TEXT symbol, delta-coded from the segment
	// base, sorted and deduplicated
	var starts []uint64
	for i := range syms {
		sp := place.sectOf(sectByName[syms[i].Section])
		if place.segs[sp.seg].name == "__TEXT" {
			starts = append(starts, addrs[i])
		}
	}
	slices.Sort(starts)
	starts = slices.Compact(starts)

	fstarts := []byte{}
	prev := uint64(machoVMAddr)
	for _, a := range starts {
		fstarts = append(fstarts, uleb(a-prev)...)
		prev = a
	}
	fstarts = append(fstarts, 0)

	// the symbol table: __mh_execute_header first (dyld looks it up),
	// then the symbols in input order; the string table opens with the
	// conventional " \0" pair so the first name starts at index 2
	strings := append([]byte{' ', 0}, "__mh_execute_header\x00"...)
	nlists := machoNlist(2, 0x0f, 1, 0x0010, machoVMAddr)
	for i := range syms {
		s := &syms[i]
		sp := place.sectOf(sectByName[s.Section])

		typ := uint8(0x0e) // N_SECT
		if s.Global {
			typ |= 0x01 // N_EXT
		}

		strings = append(strings, s.Name...)
		strings = append(strings, 0)
		nlists = append(
			nlists,
			machoNlist(uint32(len(strings)-len(s.Name)-1),
				typ, sectNum[sp.idx], 0, addrs[i])...,
		)
	}

	// regions, each padded to 8
	le.data = machoEmptyFixups()
	le.trieAt = uint32(len(le.data))
	le.data = append(le.data, trie...)
	le.data = append(le.data, make([]byte, pad8(len(le.data))-len(le.data))...)

	le.fstartAt = uint32(len(le.data))
	le.data = append(le.data, fstarts...)
	le.data = append(le.data, make([]byte, pad8(len(le.data))-len(le.data))...)

	le.symAt = uint32(len(le.data))
	le.data = append(le.data, nlists...)

	le.strAt = uint32(len(le.data))
	le.data = append(le.data, strings...)
	le.strSize = uint32(pad8(len(strings)))
	le.data = append(le.data, make([]byte, int(le.strSize)-len(strings))...)

	le.nsyms = uint32(1 + len(syms))
	le.leOff = place.total
	le.sigOff = le.leOff + uint64(len(le.data))
	le.sigLen = machoSigLen(le.sigOff)
	le.leSize = uint64(len(le.data)) + le.sigLen
	le.execSeg = uint32(place.segs[0].vmsize)

	return le
}

// machoEmptyFixups - the 56-byte chained-fixups header of an image with no
// rebase or bind targets (the constant words of the reference ld file).
func machoEmptyFixups() []byte {
	out := make([]byte, 0, 56)
	for _, v := range []uint32{0, 0x20, 0x30, 0x30, 0, 1, 0, 0, 3, 0, 0, 0, 0, 0} {
		out = binary.LittleEndian.AppendUint32(out, v)
	}

	return out
}

// machoNlist - one 16-byte nlist_64 entry.
func machoNlist(strx uint32, typ uint8, sect uint32, desc uint16, value uint64) []byte {
	out := binary.LittleEndian.AppendUint32(nil, strx)
	out = append(out, typ, uint8(sect))
	out = binary.LittleEndian.AppendUint16(out, desc)
	return binary.LittleEndian.AppendUint64(out, value)
}

// uleb - the ULEB128 encoding of v.
func uleb(v uint64) []byte {
	var out []byte
	for v >= 0x80 {
		out = append(out, byte(v)|0x80)
		v >>= 7
	}

	return append(out, byte(v))
}

// pad8 - v rounded up to a multiple of 8.
func pad8(v int) int {
	return (v + 7) &^ 7
}

// machoSigLen - the length of the ad-hoc CodeDirectory superblob for an
// image whose signature starts at sigOff (see macho_sign.go).
func machoSigLen(sigOff uint64) uint64 {
	const ident = 9 // len("assembly\x00")
	n := (sigOff + machoHashPage - 1) / machoHashPage
	return 12 + 8 + 88 + ident + 32*n
}
