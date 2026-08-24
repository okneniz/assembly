package elf

import "strconv"

// Reloc is an Elf{32,64}_Rel{,a} entry. r_info packs the symbol index and the
// type: for ELF64, sym = r_info >> 32 and type = the low 32 bits; for ELF32,
// sym = r_info >> 8 and type = the low 8 bits (the layouts are DIFFERENT
// between the classes).
//
// Elf64_Rela (24 bytes): r_offset(8) r_info(8) r_addend(8)
// Elf64_Rel  (16 bytes):  r_offset(8) r_info(8)
// Elf32_Rela (12 bytes):  r_offset(4) r_info(4) r_addend(4)
// Elf32_Rel  (8 bytes):   r_offset(4) r_info(4).
type Reloc struct {
	Off      uint64 // r_offset - the place of application (vaddr or section offset)
	SymIndex uint32 // index of the symbol in the linked symtab (0 for local ones)
	Type     uint32 // machine-dependent type
	Addend   int64  // RELA only (in REL it is embedded at the place of application)
	IsRela   bool
}

const (
	rela64Size = 24
	rel64Size  = 16
	rela32Size = 12
	rel32Size  = 8
)

// Relocations returns the relocations applied to section s: the entries of all
// SHT_RELA/SHT_REL sections whose sh_info points to s.
func (f *File) Relocations(s *Section) ([]Reloc, error) {
	var out []Reloc
	for _, rs := range f.sections {
		if rs.Info != uint32(s.index) {
			continue
		}

		var isRela bool
		var minEnt int
		switch rs.Type {
		case SHT_RELA:
			isRela, minEnt = true, rela32Size
			if f.class == CLASS64 {
				minEnt = rela64Size
			}
		case SHT_REL:
			isRela, minEnt = false, rel32Size
			if f.class == CLASS64 {
				minEnt = rel64Size
			}
		default:
			continue
		}

		if rs.Link >= uint32(len(f.sections)) {
			return nil, errf(
				"relocation section %s: sh_link=%d is outside the section table",
				rs.Name,
				rs.Link,
			)
		}

		data, err := rs.Data()
		if err != nil {
			return nil, err
		}

		entsize := rs.Entsize
		if entsize == 0 {
			entsize = uint64(minEnt)
		}

		if entsize < uint64(minEnt) {
			return nil, errf(
				"relocation section %s: entsize=%d is less than the minimum %d",
				rs.Name,
				entsize,
				minEnt,
			)
		}

		for off := 0; off+minEnt <= len(data); off += int(entsize) {
			var r Reloc
			r.IsRela = isRela
			if f.class == CLASS64 {
				r.Off = u64(data[off:], f.order)
				info := u64(data[off+8:], f.order)
				r.SymIndex = uint32(info >> 32)
				r.Type = uint32(info)
				if isRela {
					r.Addend = int64(u64(data[off+16:], f.order))
				}
			} else {
				r.Off = uint64(u32(data[off:], f.order))
				info := u32(data[off+4:], f.order)
				r.SymIndex = info >> 8
				r.Type = info & 0xff
				if isRela {
					r.Addend = int64(int32(u32(data[off+8:], f.order)))
				}
			}

			out = append(out, r)
		}
	}

	return out, nil
}

// RelocName returns the name of a relocation type for the file's machine (or hex).
func (f *File) RelocName(t uint32) string {
	switch f.hdr.Machine {
	case EM_AARCH64:
		return RelocAarch64(t).String()
	case EM_RISCV:
		return RelocRiscv(t).String()
	default:
		return strconv.FormatUint(uint64(t), 10)
	}
}
