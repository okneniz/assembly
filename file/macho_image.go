package file

// MachOImage - a native arm64 MH_EXECUTE assembled from arbitrary sections
// and symbols: the universal form of WriteMachO (which builds one and
// delegates). The image owns its whole layout - segment grouping, page
// alignment, load commands, linkedit tables, the ad-hoc signature - so no
// caller ever computes an address: the exported placement functions
// (MachoPlaceSections, MachoPlaceStreams) run the same engine, and the
// assembler layers resolve label references through them before any bytes
// reach the writer.
//
// Sections declare their segment: __TEXT (r-x, the executable image plus
// any read-only data placed there), __DATA and custom segments (rw-). A
// section is raw data or a NOBITS reserve (__bss: memory the kernel zero
// fills, no file bytes). Symbols are name + section + offset; global ones
// become external nlists and exports-trie terminals, the rest local symtab
// entries. The writer itself adds the __mh_execute_header export dyld
// expects. The entry point is a symbol name, at any offset - there is no
// template limit. Every symbol of a __TEXT section is also a function
// start (FUNCTION_STARTS is generated from them).

import (
	"errors"
	"fmt"
)

// MachOSection - one input section: raw data, or a NOBITS reserve (a bss:
// memory without file bytes, contributing to vmsize only). Section names
// are unique across the image (symbols reference sections by bare name).
// Align is bytes, a power of two; 0 means 1.
type MachOSection struct {
	Segment string
	Name    string
	Data    []byte
	Nobits  int
	Align   int
}

// machoSectMem - the memory size of a section: its data or its reserve.
func machoSectMem(s MachOSection) uint64 {
	if s.Nobits > 0 {
		return uint64(s.Nobits)
	}

	return uint64(len(s.Data))
}

// MachOSym - one symbol: a name at an offset inside a section. Global
// symbols are exported (external nlist + exports trie); locals stay in the
// symbol table only.
type MachOSym struct {
	Name    string
	Section string
	Off     uint64
	Global  bool
}

// MachOImage - the assembled image: built only by NewMachOImage (all
// validation, placement, and linkedit generation happen there), emitted
// with Bytes.
type MachOImage struct {
	sections []MachOSection
	syms     []MachOSym
	addrs    []uint64 // resolved symbol addresses, parallel to syms
	entry    int      // index into syms
	entryOff uint64   // the entry point as a file offset (LC_MAIN)
	place    machoPlacement
	le       machoLinkedit
}

// NewMachOImage validates the parts and freezes their placement. entry
// must name a symbol of a data section inside __TEXT (dyld calls it with C
// main registers; the program must not return - exit via a syscall).
func NewMachOImage(
	sections []MachOSection,
	syms []MachOSym,
	entry string,
) (*MachOImage, error) {
	if len(sections) == 0 {
		return nil, errors.New("macho: no sections")
	}

	sectIdx := make(map[string]int, len(sections))
	for i := range sections {
		s := &sections[i]

		if s.Segment == "" {
			return nil, fmt.Errorf("macho: section %q: empty segment", s.Name)
		}

		if s.Name == "" {
			return nil, fmt.Errorf("macho: section %d of %s: empty name", i, s.Segment)
		}

		if s.Segment == "__LINKEDIT" || s.Segment == "__PAGEZERO" {
			return nil, fmt.Errorf(
				"macho: section %q: segment %s is synthesized by the writer",
				s.Name,
				s.Segment,
			)
		}

		if s.Nobits > 0 && len(s.Data) > 0 {
			return nil, fmt.Errorf(
				"macho: section %s,%s: both data and a nobits reserve",
				s.Segment,
				s.Name,
			)
		}

		if s.Nobits < 0 {
			return nil, fmt.Errorf(
				"macho: section %s,%s: negative nobits reserve",
				s.Segment,
				s.Name,
			)
		}

		if s.Align < 0 || s.Align&(s.Align-1) != 0 {
			return nil, fmt.Errorf(
				"macho: section %s,%s: alignment %d is not a power of two",
				s.Segment,
				s.Name,
				s.Align,
			)
		}

		if _, dup := sectIdx[s.Name]; dup {
			return nil, fmt.Errorf("macho: duplicate section name %q", s.Name)
		}

		sectIdx[s.Name] = i
	}

	place, perr := placeMachO(sections)
	if perr != nil {
		return nil, perr
	}

	m := &MachOImage{
		sections: sections,
		syms:     syms,
		entry:    -1,
		place:    place,
	}

	m.addrs = make([]uint64, len(syms))
	for i := range syms {
		s := &syms[i]

		if s.Name == "" {
			return nil, fmt.Errorf("macho: symbol %d: empty name", i)
		}

		if s.Name == "__mh_execute_header" {
			return nil, errors.New(
				"macho: __mh_execute_header is synthesized by the writer",
			)
		}

		idx, ok := sectIdx[s.Section]
		if !ok {
			return nil, fmt.Errorf("macho: symbol %q: no section %q", s.Name, s.Section)
		}

		sp := place.sectOf(idx)
		if s.Off > sp.size {
			return nil, fmt.Errorf(
				"macho: symbol %q: offset %#x is outside section %q (%d bytes)",
				s.Name,
				s.Off,
				s.Section,
				sp.size,
			)
		}

		m.addrs[i] = sp.addr + s.Off

		if s.Name == entry {
			m.entry = i
		}
	}

	if entry == "" {
		return nil, errors.New("macho: no entry symbol")
	}

	if m.entry < 0 {
		return nil, fmt.Errorf("macho: entry symbol %q is not defined", entry)
	}

	es := &syms[m.entry]
	esp := place.sectOf(sectIdx[es.Section])
	if place.segs[esp.seg].name != "__TEXT" || esp.nobits {
		return nil, fmt.Errorf(
			"macho: entry symbol %q must sit in a __TEXT data section",
			entry,
		)
	}

	if es.Off >= esp.size {
		return nil, fmt.Errorf(
			"macho: entry symbol %q: offset %#x is at the end of %q",
			entry,
			es.Off,
			es.Section,
		)
	}

	m.entryOff = uint64(esp.offset) + es.Off
	m.le = newMachOLinkedit(&place, sectIdx, syms, m.addrs)

	if m.le.sigOff+m.le.sigLen >= machoMaxU32 {
		return nil, fmt.Errorf(
			"macho: image does not fit the 32-bit file offsets (%d bytes)",
			m.le.sigOff+m.le.sigLen,
		)
	}

	return m, nil
}

// Bytes emits the whole executable: header, load commands, section data,
// the generated linkedit tables, and the ad-hoc code signature.
func (m *MachOImage) Bytes() []byte {
	return emitMachO(m)
}

// machoMaxU32 - Mach-O file offsets and sizes are uint32; an image at the
// edge is rejected in the constructor, so the emitters cast freely.
const machoMaxU32 = uint64(1) << 32
