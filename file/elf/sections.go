package elf

// Section is a full Elf{32,64}_Shdr section header with normalized fields and
// lazy access to the data through the file's single buffer.
type Section struct {
	Name      string
	Type      SectionType
	Flags     SectionFlag
	Addr      uint64 // sh_addr - virtual address (0 for non-ALLOC)
	Off       uint64 // sh_offset - offset in the file (0 for SHT_NOBITS)
	Size      uint64 // sh_size
	Link      uint32 // sh_link - index of a linked section (strtab for symtab and the like)
	Info      uint32 // sh_info - index of the first global/target section
	Addralign uint64
	Entsize   uint64

	index int // index in the section table
	f     *File
	data  []byte // Data() cache
}

func NewSection(
	name string,
	type_ SectionType,
	flags SectionFlag,
	addr uint64,
	off uint64,
	size uint64,
	link uint32,
	info uint32,
	addralign uint64,
	entsize uint64,
	index int,
	f *File,
) *Section {
	return &Section{
		Name:      name,
		Type:      type_,
		Flags:     flags,
		Addr:      addr,
		Off:       off,
		Size:      size,
		Link:      link,
		Info:      info,
		Addralign: addralign,
		Entsize:   entsize,
		index:     index,
		f:         f,
	}
}

// Data returns the section contents: read lazily from the buffer and cached.
// For SHT_NOBITS it returns zero bytes (there is no data in the file).
func (s *Section) Data() ([]byte, error) {
	if s.data != nil {
		return s.data, nil
	}

	if s.Type == SHT_NOBITS || s.Size == 0 {
		s.data = []byte{}
		return s.data, nil
	}

	data, err := readBytes(s.f.buf, int(s.Off), int(s.Size))
	if err != nil {
		return nil, err
	}

	s.data = data
	return data, nil
}

// Field offsets of Elf64_Shdr (size 64):
//
//	sh_name(4)@0 sh_type(4)@4 sh_flags(8)@8 sh_addr(8)@16 sh_offset(8)@24
//	sh_size(8)@32 sh_link(4)@40 sh_info(4)@44 sh_addralign(8)@48 sh_entsize(8)@56
//
// Elf32_Shdr (size 40):
//
//	sh_name(4)@0 sh_type(4)@4 sh_flags(4)@8 sh_addr(4)@12 sh_offset(4)@16
//	sh_size(4)@20 sh_link(4)@24 sh_info(4)@28 sh_addralign(4)@32 sh_entsize(4)@36
const (
	shdr64Size  = 64
	shdr32Size  = 40
	shdr64Flags = 8
	shdr64Addr  = 16
	shdr64Off   = 24
	shdr64SizeF = 32
	shdr64Link  = 40
	shdr64Info  = 44
	shdr64Align = 48
	shdr64Entsz = 56
	shdr32Flags = 8
	shdr32Addr  = 12
	shdr32Off   = 16
	shdr32SizeF = 20
	shdr32Link  = 24
	shdr32Info  = 28
	shdr32Align = 32
	shdr32Entsz = 36
)

// Sections returns all sections of the file (the headers are parsed at Open;
// names resolved through shstrtab).
func (f *File) Sections() []*Section {
	return f.sections
}

// Section returns the section by name (e.g. ".text"), or nil.
func (f *File) Section(name string) *Section {
	for _, s := range f.sections {
		if s.Name == name {
			return s
		}
	}

	return nil
}

// parseSections reads the section header table and assigns the sections names
// from the shstrndx string table (if there is one).
func (f *File) parseSections() error {
	if f.hdr.Shnum == 0 || f.hdr.Shoff == 0 {
		return nil
	}

	size := shdr32Size
	if f.class == CLASS64 {
		size = shdr64Size
	}

	if int(f.hdr.Shentsize) < size {
		return errf("e_shentsize=%d is less than the minimum %d", f.hdr.Shentsize, size)
	}

	if f.hdr.Shstrndx != 0 && f.hdr.Shstrndx >= f.hdr.Shnum {
		return errf("e_shstrndx=%d is outside the section table (%d)", f.hdr.Shstrndx, f.hdr.Shnum)
	}

	type shdrRaw struct {
		nameOff, typ uint32
		flags        uint64
		addr, off    uint64
		size         uint64
		link, info   uint32
		align, entsz uint64
	}
	raws := make([]shdrRaw, f.hdr.Shnum)
	for i := range f.hdr.Shnum {
		base := int(f.hdr.Shoff) + i*int(f.hdr.Shentsize)
		r := &raws[i]
		if f.class == CLASS64 {
			var e error
			if r.nameOff, e = readU32At(f.buf, f.order, base); e != nil {
				return e
			}

			if r.typ, e = readU32At(f.buf, f.order, base+4); e != nil {
				return e
			}

			if r.flags, e = readU64At(f.buf, f.order, base+shdr64Flags); e != nil {
				return e
			}

			if r.addr, e = readU64At(f.buf, f.order, base+shdr64Addr); e != nil {
				return e
			}

			if r.off, e = readU64At(f.buf, f.order, base+shdr64Off); e != nil {
				return e
			}

			if r.size, e = readU64At(f.buf, f.order, base+shdr64SizeF); e != nil {
				return e
			}

			if r.link, e = readU32At(f.buf, f.order, base+shdr64Link); e != nil {
				return e
			}

			if r.info, e = readU32At(f.buf, f.order, base+shdr64Info); e != nil {
				return e
			}

			if r.align, e = readU64At(f.buf, f.order, base+shdr64Align); e != nil {
				return e
			}

			if r.entsz, e = readU64At(f.buf, f.order, base+shdr64Entsz); e != nil {
				return e
			}

			continue
		}

		var e error
		var v uint32
		if r.nameOff, e = readU32At(f.buf, f.order, base); e != nil {
			return e
		}

		if r.typ, e = readU32At(f.buf, f.order, base+4); e != nil {
			return e
		}

		if v, e = readU32At(f.buf, f.order, base+shdr32Flags); e != nil {
			return e
		}

		r.flags = uint64(v)
		if v, e = readU32At(f.buf, f.order, base+shdr32Addr); e != nil {
			return e
		}

		r.addr = uint64(v)
		if v, e = readU32At(f.buf, f.order, base+shdr32Off); e != nil {
			return e
		}

		r.off = uint64(v)
		if v, e = readU32At(f.buf, f.order, base+shdr32SizeF); e != nil {
			return e
		}

		r.size = uint64(v)
		if r.link, e = readU32At(f.buf, f.order, base+shdr32Link); e != nil {
			return e
		}

		if r.info, e = readU32At(f.buf, f.order, base+shdr32Info); e != nil {
			return e
		}

		if v, e = readU32At(f.buf, f.order, base+shdr32Align); e != nil {
			return e
		}

		r.align = uint64(v)
		if v, e = readU32At(f.buf, f.order, base+shdr32Entsz); e != nil {
			return e
		}

		r.entsz = uint64(v)
	}

	// The section-name string table is loaded into memory whole; names are
	// sliced from it.
	strtab := []byte{}
	if f.hdr.Shstrndx > 0 {
		sr := raws[f.hdr.Shstrndx]
		if sr.typ != uint32(SHT_STRTAB) {
			return errf("shstrndx points to %s, want SHT_STRTAB", SectionType(sr.typ))
		}

		if sr.size > 0 {
			data, e := readBytes(f.buf, int(sr.off), int(sr.size))
			if e != nil {
				return e
			}

			strtab = data
		}
	}

	f.sections = make([]*Section, f.hdr.Shnum)
	for i, r := range raws {
		name := ""
		if f.hdr.Shstrndx > 0 && r.nameOff < uint32(len(strtab)) {
			name = cstr(strtab[r.nameOff:])
		}

		f.sections[i] = NewSection(
			name,
			SectionType(r.typ),
			SectionFlag(r.flags),
			r.addr,
			r.off,
			r.size,
			r.link,
			r.info,
			r.align,
			r.entsz,
			i,
			f,
		)
	}

	return nil
}

// cstr returns the prefix of b up to the first null byte.
func cstr(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}

	return string(b)
}
