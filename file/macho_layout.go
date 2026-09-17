package file

// The placement engine of the Mach-O writer: the ONE place that knows the
// image geometry. Everything else - the exported placement policies for
// the assembler layers, the emitters, the linkedit generators - consumes
// its output, which is what keeps the library contract-free: a section
// list goes in, final addresses come out, and the same numbers place the
// bytes in the file.
//
// Layout rules (Apple Silicon hard requirements - see macho_write.go):
//   - __TEXT is first at the arm64 convention base, file offset 0,
//     carrying the header and the load commands before its sections;
//   - every segment starts file- and memory-page aligned (16K), and every
//     mapped page is fully file-backed except a trailing NOBITS reserve;
//   - inside a segment the data sections come first (input order), NOBITS
//     reserves last (input order) - like ld, so the zero-fill tail never
//     lands in the middle of file-backed pages;
//   - __LINKEDIT (the generated tables + the signature) follows the last
//     segment, page aligned.

import (
	"errors"
	"math/bits"
)

// machoSectPlace - the frozen placement of one input section.
type machoSectPlace struct {
	idx    int    // index into the input section list
	seg    int    // index into place.segs
	addr   uint64 // vmaddr
	offset uint32 // file offset of the data (nobits: no file bytes)
	align  uint32 // 2^align bytes
	size   uint64 // memory size
	nobits bool
}

// machoSegPlace - one segment and its sections.
type machoSegPlace struct {
	name     string
	vmaddr   uint64
	vmsize   uint64
	fileoff  uint64
	filesize uint64
	maxprot  uint32
	initprot uint32
	sects    []machoSectPlace
}

// machoPlacement - the geometry of a whole image.
type machoPlacement struct {
	segs    []machoSegPlace
	sect    []machoSectPlace // parallel to the input section list
	ncmds   uint32
	cmdSize uint32
	codeOff uint32 // file offset of the first __TEXT section byte
	total   uint64 // file offset where __LINKEDIT starts (page aligned)
	vmTotal uint64 // __LINKEDIT vmaddr - machoVMAddr (vm address continuation)
}

// sectOf - the placement of input section idx.
func (p *machoPlacement) sectOf(idx int) *machoSectPlace {
	return &p.sect[idx]
}

// machoFixedCmds - the size of every load command except the segments:
// chained fixups, exports trie, symtab, dyld info, dylinker, uuid, build
// version, source version, main, libSystem, function starts, data in
// code, code signature.
const machoFixedCmds = 368

// machoSpec - the placement input of one section: its identity, memory
// size, and alignment - no data. The exported policies predict layouts
// from sizes alone; placeMachO derives the specs from full sections.
type machoSpec struct {
	segment string
	name    string
	size    uint64
	nobits  bool
	align   int
}

// placeMachO groups the sections into segments and freezes the geometry.
// The input is already validated (NewMachOImage): here it only lays out.
func placeMachO(sections []MachOSection) (machoPlacement, error) {
	if len(sections) == 0 {
		return machoPlacement{}, errors.New("macho: no sections")
	}

	specs := make([]machoSpec, len(sections))
	for i := range sections {
		s := &sections[i]
		specs[i] = machoSpec{
			segment: s.Segment,
			name:    s.Name,
			size:    machoSectMem(*s),
			nobits:  s.Nobits > 0,
			align:   s.Align,
		}
	}

	return placeMachOSpecs(specs)
}

// placeMachOSpecs - the placement engine itself, over size-only specs.
func placeMachOSpecs(specs []machoSpec) (machoPlacement, error) {
	var p machoPlacement
	if len(specs) == 0 {
		return p, errors.New("macho: no sections")
	}

	// segment list: __TEXT first, then the rest in order of appearance
	segIdx := make(map[string]int)
	for i := range specs {
		s := &specs[i]

		if _, ok := segIdx[s.segment]; !ok {
			if len(p.segs) > 0 && s.segment == "__TEXT" {
				return p, errors.New(
					"macho: __TEXT sections must come before the other segments",
				)
			}

			segIdx[s.segment] = len(p.segs)
			p.segs = append(p.segs, machoSegPlace{name: s.segment})
		}
	}

	if _, ok := segIdx["__TEXT"]; !ok {
		return p, errors.New("macho: no __TEXT section")
	}

	// sections into their segments: data first, nobits reserves last
	for i := range specs {
		s := &specs[i]
		seg := &p.segs[segIdx[s.segment]]

		align := uint32(0)
		if s.align > 0 {
			align = uint32(bits.TrailingZeros64(uint64(s.align)))
		}

		at := len(seg.sects)
		if !s.nobits {
			for at > 0 && seg.sects[at-1].nobits {
				at--
			}
		}

		sp := machoSectPlace{
			idx:    i,
			seg:    segIdx[s.segment],
			align:  align,
			size:   s.size,
			nobits: s.nobits,
		}

		seg.sects = append(seg.sects, sp)
		copy(seg.sects[at+1:], seg.sects[at:])
		seg.sects[at] = sp
	}

	// command geometry: the fixed set, PAGEZERO and __LINKEDIT, one
	// command per segment, 80 bytes per section
	p.ncmds = uint32(13 + 2 + len(p.segs))
	p.cmdSize = uint32(machoFixedCmds + 72*(2+len(p.segs)) + 80*len(specs))

	// segment geometry, __TEXT first at the convention base. The file and
	// memory cursors advance independently: a segment with no file bytes
	// (a pure bss) points at file offset 0, like ld, and does not move
	// the file cursor - only its memory pages continue the address space.
	vm := uint64(machoVMAddr)
	file := uint64(0)
	for si := range p.segs {
		seg := &p.segs[si]
		seg.vmaddr = vm
		seg.fileoff = file // reset to 0 below when the segment has no bytes

		start := uint64(0)
		if seg.name == "__TEXT" {
			seg.maxprot, seg.initprot = 7, 5
			start = alignUp(uint64(32+p.cmdSize),
				uint64(specs[seg.sects[0].idx].align))
			p.codeOff = uint32(start)
		} else {
			seg.maxprot, seg.initprot = 3, 3
		}

		placeSegSections(seg, start)

		seg.filesize = roundMachoKPage(seg.fileSpan())
		seg.vmsize = roundMachoKPage(seg.memSpan())
		if seg.filesize > 0 {
			file += seg.filesize
		} else {
			seg.fileoff = 0
		}

		vm = seg.vmaddr + seg.vmsize
	}

	// the section view, parallel to the input list
	p.sect = make([]machoSectPlace, len(specs))
	for si := range p.segs {
		for _, sp := range p.segs[si].sects {
			p.sect[sp.idx] = sp
		}
	}

	p.total = file
	p.vmTotal = vm - machoVMAddr

	return p, nil
}

// placeSegSections walks the sections of one segment: every section is
// aligned inside the segment, data sections advance the file cursor too,
// nobits reserves advance memory only (they sort last, so the file-backed
// span stays contiguous). The addresses land in the segment's own sects
// slice; p.sect is rebuilt from them by the caller.
func placeSegSections(seg *machoSegPlace, start uint64) {
	cur := start
	file := start
	first := true

	for i := range seg.sects {
		sp := &seg.sects[i]

		if !first {
			cur = alignUp(cur, uint64(1)<<sp.align)
			file = alignUp(file, uint64(1)<<sp.align)
		}
		first = false

		sp.addr = seg.vmaddr + cur
		cur += sp.size

		if !sp.nobits {
			sp.offset = uint32(seg.fileoff + file)
			file += sp.size
		}
	}
}

// fileSpan - the offset just past the last file byte of the segment.
func (s *machoSegPlace) fileSpan() uint64 {
	span := uint64(0)
	for _, sp := range s.sects {
		if !sp.nobits {
			span = uint64(sp.offset) - s.fileoff + sp.size
		}
	}

	return span
}

// memSpan - the offset just past the last section byte in memory.
func (s *machoSegPlace) memSpan() uint64 {
	span := uint64(0)
	for _, sp := range s.sects {
		if end := sp.addr - s.vmaddr + sp.size; end > span {
			span = end
		}
	}

	return span
}

// alignUp rounds v up to a multiple of a (a is a power of two, or 1).
func alignUp(v, a uint64) uint64 {
	if a <= 1 {
		return v
	}

	return (v + a - 1) &^ (a - 1)
}

// roundMachoKPage rounds v up to whole 16K kernel pages.
func roundMachoKPage(v uint64) uint64 {
	return alignUp(v, machoKPage)
}
