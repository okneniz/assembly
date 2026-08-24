package macho

import "errors"

// Chained fixups (LC_DYLD_CHAINED_FIXUPS): the modern replacement for the
// classic rebase/bind streams. Layouts are from <mach-o/fixup-chains.h>.

// Pointer and import formats.
const (
	DYLD_CHAINED_IMPORT          uint32 = 1 // 4 bytes
	DYLD_CHAINED_IMPORT_ADDEND   uint32 = 2 // 8 bytes (+int32)
	DYLD_CHAINED_IMPORT_ADDEND64 uint32 = 3 // 16 bytes (+int64)

	DYLD_CHAINED_PTR_ARM64E              uint16 = 1 // stride 4
	DYLD_CHAINED_PTR_64                  uint16 = 2 // stride 4
	DYLD_CHAINED_PTR_32                  uint16 = 3
	DYLD_CHAINED_PTR_32_CACHE            uint16 = 4
	DYLD_CHAINED_PTR_32_FIRMWARE         uint16 = 5
	DYLD_CHAINED_PTR_64_OFFSET           uint16 = 6 // stride 4
	DYLD_CHAINED_PTR_ARM64E_KERNEL       uint16 = 7 // stride 4
	DYLD_CHAINED_PTR_64_KERNEL_CACHE     uint16 = 8
	DYLD_CHAINED_PTR_ARM64E_USERLAND     uint16 = 9  // stride 8
	DYLD_CHAINED_PTR_ARM64E_FIRMWARE     uint16 = 10 // stride 4
	DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE uint16 = 11 // stride 1
	DYLD_CHAINED_PTR_ARM64E_USERLAND24   uint16 = 12 // stride 8

	DYLD_CHAINED_PTR_START_NONE  uint16 = 0xffff // page with no fixups
	DYLD_CHAINED_PTR_START_MULTI uint16 = 0x8000 // multiple chains per page
	DYLD_CHAINED_PTR_START_LAST  uint16 = 0x8000 // last start in chain_starts
)

// ChainedImport is an entry of the fixups import table.
type ChainedImport struct {
	LibOrdinal int
	WeakImport bool
	Name       string
	Addend     int64 // for the ADDEND/ADDEND64 formats
}

// Fixup is a single chain element: a rebase (pointer fixup) or a bind
// (binding to an import).
type Fixup struct {
	Addr uint64 // vmaddr of the fixup site
	Bind bool   // false means rebase

	// Rebase.
	Target uint64
	High8  uint8 // for PTR_64: high 8 bits of the pointer

	// Bind.
	ImportIndex uint32
	Addend      int64

	// arm64e authenticated pointers.
	Auth      bool
	Diversity uint16
	AddrDiv   bool
	Key       uint8
}

func NewFixup(addr uint64) Fixup {
	return Fixup{Addr: addr}
}

// ChainedFixups is the parsed payload of LC_DYLD_CHAINED_FIXUPS.
type ChainedFixups struct {
	Version       uint32
	ImportsFormat uint32
	SymbolsFormat uint32
	Imports       []ChainedImport
	Fixups        []Fixup
}

func NewChainedFixups(version uint32, importsFormat uint32, symbolsFormat uint32) *ChainedFixups {
	return &ChainedFixups{
		Version:       version,
		ImportsFormat: importsFormat,
		SymbolsFormat: symbolsFormat,
	}
}

// pointerStride returns the chain stride (in bytes) for a pointer format.
func pointerStride(format uint16) int {
	switch format {
	case DYLD_CHAINED_PTR_ARM64E_USERLAND, DYLD_CHAINED_PTR_ARM64E_USERLAND24:
		return 8
	case DYLD_CHAINED_PTR_X86_64_KERNEL_CACHE:
		return 1
	}

	return 4
}

// ErrNoFixups: the file has no LC_DYLD_CHAINED_FIXUPS load command.
var ErrNoFixups = errors.New("macho: no LC_DYLD_CHAINED_FIXUPS load command")

// Fixups parses LC_DYLD_CHAINED_FIXUPS; a missing command corresponds to
// ErrNoFixups.
func (f *File) Fixups() (*ChainedFixups, error) {
	var payload []byte
	for _, lc := range f.commands {
		if ld, ok := lc.(*LinkeditData); ok && ld.Cmd() == LC_DYLD_CHAINED_FIXUPS {
			data, err := f.dyldStream(ld.DataOff, ld.DataSize)
			if err != nil {
				return nil, err
			}

			payload = data
		}
	}

	if payload == nil {
		return nil, ErrNoFixups
	}

	if len(payload) < 28 {
		return nil, errf("chained fixups: truncated header (%d bytes)", len(payload))
	}

	cx := NewChainedFixups(
		u32(payload, f.order),
		u32(payload[20:], f.order),
		u32(payload[24:], f.order),
	)
	startsOffset := u32(payload[4:], f.order)
	importsOffset := u32(payload[8:], f.order)
	symbolsOffset := u32(payload[12:], f.order)
	importsCount := u32(payload[16:], f.order)

	// Import table + string pool.
	symbols := payload[min(int(symbolsOffset), len(payload)):]
	importEntSize := 4
	switch cx.ImportsFormat {
	case DYLD_CHAINED_IMPORT_ADDEND:
		importEntSize = 8
	case DYLD_CHAINED_IMPORT_ADDEND64:
		importEntSize = 16
	case DYLD_CHAINED_IMPORT:
	default:
		return nil, errf("chained fixups: unknown imports format %d", cx.ImportsFormat)
	}

	for i := range importsCount {
		off := int(importsOffset) + int(i)*importEntSize
		if off+importEntSize > len(payload) {
			return nil, errf("chained fixups: import %d out of data", i)
		}

		var im ChainedImport
		if cx.ImportsFormat == DYLD_CHAINED_IMPORT_ADDEND64 {
			w := readU64(payload[off:], f.order)
			im.LibOrdinal = int(w & 0xffff)
			im.WeakImport = w>>16&1 != 0
			nameOff := uint32(w >> 32)
			im.Addend = int64(readU64(payload[off+8:], f.order))
			im.Name = cstrAt(symbols, nameOff)
		} else {
			w := u32(payload[off:], f.order)
			im.LibOrdinal = int(w & 0xff)
			im.WeakImport = w>>8&1 != 0
			nameOff := w >> 9
			if cx.ImportsFormat == DYLD_CHAINED_IMPORT_ADDEND {
				im.Addend = int64(int32(u32(payload[off+4:], f.order)))
			}

			im.Name = cstrAt(symbols, nameOff)
		}

		cx.Imports = append(cx.Imports, im)
	}

	// starts_in_image: seg_count + offsets of starts_in_segment.
	if int(startsOffset)+4 > len(payload) {
		return nil, errf("chained fixups: truncated starts_in_image")
	}

	segCount := u32(payload[startsOffset:], f.order)
	for i := range segCount {
		infoOffAt := int(startsOffset) + 4 + int(i)*4
		if infoOffAt+4 > len(payload) {
			return nil, errf("chained fixups: seg_info_offset[%d] out of data", i)
		}

		segInfoOff := u32(payload[infoOffAt:], f.order)
		if segInfoOff == 0 {
			continue
		}

		if int(i) >= len(f.segments) {
			break
		}

		// seg_info_offset is relative to the start of starts_in_image.
		if err := f.walkSegmentChains(
			payload,
			int(startsOffset)+int(segInfoOff),
			f.segments[i],
			cx,
		); err != nil {
			return nil, err
		}
	}

	return cx, nil
}

// walkSegmentChains walks the chains of a single segment.
//
// starts_in_segment: size(4) page_size(2) pointer_format(2) segment_offset(8,
// vmaddr) max_valid_pointer(4) page_count(2) page_start[page_count].
func (f *File) walkSegmentChains(payload []byte, off int, seg *Segment, cx *ChainedFixups) error {
	base := off
	if base+24 > len(payload) {
		return errf("chained fixups: truncated starts_in_segment")
	}

	pageSize := readU16(payload[base+4:], f.order)
	format := readU16(payload[base+6:], f.order)
	segOffset := readU64(payload[base+8:], f.order) // vmaddr of the segment
	pageCount := readU16(payload[base+20:], f.order)

	// The segment in the payload is identified by vmaddr: match it against
	// the file's segments (for fileoff we use the segment with the same vmaddr).
	fileSeg := seg
	if segOffset != seg.Vmaddr {
		for _, s := range f.segments {
			if s.Vmaddr == segOffset {
				fileSeg = s
				break
			}
		}
	}

	stride := pointerStride(format)
	for page := range pageCount {
		startAt := base + 22 + int(page)*2
		if startAt+2 > len(payload) {
			return errf("chained fixups: page_start[%d] out of data", page)
		}

		start := readU16(payload[startAt:], f.order)
		if start == DYLD_CHAINED_PTR_START_NONE {
			continue
		}

		if start&DYLD_CHAINED_PTR_START_MULTI != 0 {
			// Multiple chains per page (32-bit formats): a page_start with
			// the MULTI bit is an INDEX into the chain_starts array, which
			// lies right after page_start[]; the last element is marked
			// with bit 0x8000.
			overflowBase := base + 22 + int(pageCount)*2
			listAt := overflowBase + int(start&0x7fff)*2
			for {
				if listAt+2 > len(payload) {
					return errf("chained fixups: truncated chain_starts")
				}

				s2 := readU16(payload[listAt:], f.order)
				if err := f.walkChain(
					fileSeg,
					format,
					stride,
					uint64(page)*uint64(pageSize),
					uint64(s2&0x7fff),
					cx,
				); err != nil {
					return err
				}

				listAt += 2
				if s2&DYLD_CHAINED_PTR_START_LAST != 0 {
					break
				}
			}

			continue
		}

		if err := f.walkChain(
			fileSeg,
			format,
			stride,
			uint64(page)*uint64(pageSize),
			uint64(start),
			cx,
		); err != nil {
			return err
		}
	}

	return nil
}

// walkChain walks a single chain of fixups from a start offset within the page.
func (f *File) walkChain(
	seg *Segment,
	format uint16,
	stride int,
	pageBase, start uint64,
	cx *ChainedFixups,
) error {
	is32 := format == DYLD_CHAINED_PTR_32 || format == DYLD_CHAINED_PTR_32_CACHE ||
		format == DYLD_CHAINED_PTR_32_FIRMWARE
	entrySize := uint64(8)
	if is32 {
		entrySize = 4
	}

	offset := start
	for {
		raw, err := readBytes(f.buf, f.base+int(seg.Fileoff)+int(pageBase+offset), int(entrySize))
		if err != nil {
			return err
		}

		var w uint64
		if is32 {
			w = uint64(u32(raw, f.order))
		} else {
			w = readU64(raw, f.order)
		}

		fx := NewFixup(seg.Vmaddr + pageBase + offset)

		var next uint64
		if is32 {
			next = w >> 27 & 0x1f
			if w&1 != 0 { // bind
				fx.Bind = true
				fx.ImportIndex = uint32(w>>6) & 0xfffff
				fx.Addend = int64(w>>1) & 0x3f
			} else {
				fx.Target = w >> 6 & 0x3ffffff
			}
		} else {
			switch format {
			case DYLD_CHAINED_PTR_ARM64E, DYLD_CHAINED_PTR_ARM64E_KERNEL,
				DYLD_CHAINED_PTR_ARM64E_USERLAND, DYLD_CHAINED_PTR_ARM64E_USERLAND24,
				DYLD_CHAINED_PTR_ARM64E_FIRMWARE:
				next = w >> 51 & 0x7ff
				auth := w>>63&1 != 0
				bind := w>>62&1 != 0
				fx.Auth = auth
				switch {
				case bind:
					fx.Bind = true
					if format == DYLD_CHAINED_PTR_ARM64E_USERLAND24 {
						fx.ImportIndex = uint32(w & 0xffffff)
					} else {
						fx.ImportIndex = uint32(w & 0xffff)
					}

					if !auth {
						fx.Addend = sextN(w>>32&0x7ffff, 19)
					} else {
						fx.Diversity = uint16(w >> 32 & 0xffff)
						fx.AddrDiv = w>>48&1 != 0
						fx.Key = uint8(w >> 49 & 0x3)
					}
				case auth:
					fx.Target = w & 0xffffffff
					fx.Diversity = uint16(w >> 32 & 0xffff)
					fx.AddrDiv = w>>48&1 != 0
					fx.Key = uint8(w >> 49 & 0x3)
				default:
					fx.Target = w & 0x7ffffffffff
					fx.High8 = uint8(w >> 43 & 0xff)
				}

			default: // DYLD_CHAINED_PTR_64, _64_OFFSET, kernel caches
				// dyld_chained_ptr_64_rebase:  target:36 high8:8 reserved:7 next:12 bind:1
				// dyld_chained_ptr_64_bind:    ordinal:24 addend:8 reserved:19 next:12 bind:1
				next = w >> 51 & 0xfff
				if w>>63&1 != 0 {
					fx.Bind = true
					fx.ImportIndex = uint32(w & 0xffffff)
					fx.Addend = sextN(w>>24&0xff, 8)
				} else {
					fx.Target = w & 0xfffffffff
					fx.High8 = uint8(w >> 36 & 0xff)
				}
			}
		}

		if fx.Bind {
			if int(fx.ImportIndex) < len(cx.Imports) {
				fx.Addend += cx.Imports[fx.ImportIndex].Addend
			}
		}

		cx.Fixups = append(cx.Fixups, fx)

		if next == 0 {
			return nil
		}

		offset += next * uint64(stride)
		if pageBase+offset >= seg.Vmsize && pageBase+offset >= seg.Filesize {
			return nil // guard against running past the segment
		}
	}
}

// sextN sign-extends an n-bit value.
func sextN(v uint64, bits uint) int64 {
	if v&(1<<(bits-1)) != 0 {
		return int64(v | ^uint64(0)<<bits)
	}

	return int64(v)
}

// cstrAt reads a C string from a pool by offset.
func cstrAt(pool []byte, off uint32) string {
	if int(off) >= len(pool) {
		return ""
	}

	return cstr(pool[off:])
}
