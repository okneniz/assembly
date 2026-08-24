package elf

// Prog is an Elf{32,64}_Phdr program header table entry with normalized
// fields (32-bit ones widened to uint64).
//
// Elf64_Phdr: p_type(4) p_flags(4) p_offset(8) p_vaddr(8) p_paddr(8)
//
//	p_filesz(8) p_memsz(8) p_align(8) - 56 bytes, flags as the SECOND field.
//
// Elf32_Phdr: p_type(4) p_offset(4) p_vaddr(4) p_paddr(4)
//
//	p_filesz(4) p_memsz(4) p_flags(4) p_align(4) - 32 bytes, flags AFTER
//	p_memsz - the classic layout trap.
type Prog struct {
	Type   ProgType
	Flags  ProgFlag
	Off    uint64 // p_offset - the segment's offset in the file
	Vaddr  uint64
	Paddr  uint64
	Filesz uint64
	Memsz  uint64
	Align  uint64
}

const (
	phdr64Size  = 56
	phdr32Size  = 32
	phdr64Flags = 4
	phdr32Flags = 24
)

// ProgramHeaders returns all program headers (parsed at Open).
func (f *File) ProgramHeaders() []Prog {
	return f.progs
}

// parseProgs reads the Elf{32,64}_Phdr array at e_phoff/e_phnum with a stride
// of e_phentsize (the stride may exceed the minimal size - the extra is skipped).
func (f *File) parseProgs() error {
	if f.hdr.Phnum == 0 || f.hdr.Phoff == 0 {
		return nil
	}

	size := phdr32Size
	flagsOff := phdr32Flags
	if f.class == CLASS64 {
		size = phdr64Size
		flagsOff = phdr64Flags
	}

	if int(f.hdr.Phentsize) < size {
		return errf("e_phentsize=%d is less than the minimum %d", f.hdr.Phentsize, size)
	}

	f.progs = make([]Prog, 0, f.hdr.Phnum)
	for i := range f.hdr.Phnum {
		base := int(f.hdr.Phoff) + i*int(f.hdr.Phentsize)
		if f.class == CLASS64 {
			p := Prog{}
			var raw uint32
			var err error
			if raw, err = readU32At(f.buf, f.order, base); err != nil {
				return err
			}

			p.Type = ProgType(raw)
			if raw, err = readU32At(f.buf, f.order, base+flagsOff); err != nil {
				return err
			}

			p.Flags = ProgFlag(raw)
			if p.Off, err = readU64At(f.buf, f.order, base+8); err != nil {
				return err
			}

			if p.Vaddr, err = readU64At(f.buf, f.order, base+16); err != nil {
				return err
			}

			if p.Paddr, err = readU64At(f.buf, f.order, base+24); err != nil {
				return err
			}

			if p.Filesz, err = readU64At(f.buf, f.order, base+32); err != nil {
				return err
			}

			if p.Memsz, err = readU64At(f.buf, f.order, base+40); err != nil {
				return err
			}

			if p.Align, err = readU64At(f.buf, f.order, base+48); err != nil {
				return err
			}

			f.progs = append(f.progs, p)
			continue
		}

		p := Prog{}
		var raw, v uint32
		var err error
		if raw, err = readU32At(f.buf, f.order, base); err != nil {
			return err
		}

		p.Type = ProgType(raw)
		if v, err = readU32At(f.buf, f.order, base+4); err != nil {
			return err
		}

		p.Off = uint64(v)
		if v, err = readU32At(f.buf, f.order, base+8); err != nil {
			return err
		}

		p.Vaddr = uint64(v)
		if v, err = readU32At(f.buf, f.order, base+12); err != nil {
			return err
		}

		p.Paddr = uint64(v)
		if v, err = readU32At(f.buf, f.order, base+16); err != nil {
			return err
		}

		p.Filesz = uint64(v)
		if v, err = readU32At(f.buf, f.order, base+20); err != nil {
			return err
		}

		p.Memsz = uint64(v)
		if raw, err = readU32At(f.buf, f.order, base+flagsOff); err != nil {
			return err
		}

		p.Flags = ProgFlag(raw)
		if v, err = readU32At(f.buf, f.order, base+28); err != nil {
			return err
		}

		p.Align = uint64(v)
		f.progs = append(f.progs, p)
	}

	return nil
}
