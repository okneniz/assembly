package macho

import "strconv"

// Reloc is relocation_info (8 bytes): r_address(4) + a packed word
// r_symbolnum:24 r_pcrel:1 r_length:2 r_extern:1 r_type:4.
//
// Scatter variant (bit 0x80000000 in r_address, an obsolete 32-bit
// mechanism): r_scattered:1 r_address:24 r_type:4 r_length:2 r_pcrel:1 + r_value.
type Reloc struct {
	Addr      uint32 // offset within the section
	SymIndex  uint32 // symbol (Extern) or section index
	Pcrel     bool
	Length    uint8 // log2 of size: 0=1b, 1=2b, 2=4b, 3=8b
	Extern    bool
	Type      RelocType
	Scattered bool
	Value     uint32 // scattered only
}

func NewReloc(
	addr uint32,
	symIndex uint32,
	pcrel bool,
	length uint8,
	extern bool,
	type_ RelocType,
	scattered bool,
	value uint32,
) Reloc {
	return Reloc{
		Addr:      addr,
		SymIndex:  symIndex,
		Pcrel:     pcrel,
		Length:    length,
		Extern:    extern,
		Type:      type_,
		Scattered: scattered,
		Value:     value,
	}
}

// Addend returns the addend of a paired ARM64_RELOC_ADDEND relocation: by
// convention the ADDEND entry is followed by BRANCH26/PAGE21/PAGEOFF12,
// whose addend is stored in the pair's SymIndex.
func (r Reloc) Addend(next Reloc) int32 {
	if r.Type != ARM64_RELOC_ADDEND {
		return 0
	}

	return int32(next.SymIndex)
}

// RelocName returns the relocation type name for the file's architecture.
func (f *File) RelocName(t RelocType) string {
	switch f.hdr.CpuType {
	case CPU_TYPE_ARM64, CPU_TYPE_ARM64_32:
		return arm64RelocName(t)
	case CPU_TYPE_X86_64, CPU_TYPE_X86:
		return strconv.FormatUint(uint64(t), 10)
	}

	return strconv.FormatUint(uint64(t), 10)
}

// Relocations returns a section's relocation table (section.nreloc entries
// of 8 bytes from section.reloff; offsets are from the start of the Mach-O
// object).
func (f *File) Relocations(s *Section) ([]Reloc, error) {
	if s.Nreloc == 0 {
		return nil, nil
	}

	data, err := readBytes(f.buf, f.base+int(s.Reloff), int(s.Nreloc)*8)
	if err != nil {
		return nil, err
	}

	out := make([]Reloc, 0, s.Nreloc)
	for i := range s.Nreloc {
		w0 := u32(data[i*8:], f.order)
		w1 := u32(data[i*8+4:], f.order)

		if w0&R_SCATTERED != 0 {
			out = append(
				out,
				NewReloc(
					w0&0x00ffffff,
					0,
					w0>>30&1 != 0,
					uint8(w0>>28&0x3),
					false,
					RelocType(w0>>24&0xf),
					true,
					w1,
				),
			)
			continue
		}

		out = append(
			out,
			NewReloc(
				w0,
				w1&0x00ffffff,
				w1>>24&1 != 0,
				uint8(w1>>25&0x3),
				w1>>27&1 != 0,
				RelocType(w1>>28),
				false,
				0,
			),
		)
	}

	return out, nil
}

// IndirectSymbols returns the indirect symbol table (LC_DYSYMTAB): symtab
// indexes for pointer sections (__got, __la_symbol_ptr, __stub...).
func (f *File) IndirectSymbols() ([]uint32, error) {
	var cmd *Dysymtab
	for _, lc := range f.commands {
		if d, ok := lc.(*Dysymtab); ok {
			cmd = d
			break
		}
	}

	if cmd == nil || cmd.NIndirectSyms == 0 {
		return nil, nil
	}

	data, err := readBytes(f.buf, f.base+int(cmd.IndirectSymOff), int(cmd.NIndirectSyms)*4)
	if err != nil {
		return nil, err
	}

	out := make([]uint32, 0, cmd.NIndirectSyms)
	for i := range cmd.NIndirectSyms {
		out = append(out, u32(data[i*4:], f.order))
	}

	return out, nil
}
