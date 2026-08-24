package macho

// Segment is segment_command_64 or segment_command with nested sections.
//
// segment_command_64 (72 bytes):
//
//	cmd(4) cmdsize(4) segname(16) vmaddr(8) vmsize(8) fileoff(8) filesize(8)
//	maxprot(4) initprot(4) nsects(4) flags(4); section_64 sections, 80 bytes each
//
// segment_command (56 bytes):
//
//	cmd(4) cmdsize(4) segname(16) vmaddr(4) vmsize(4) fileoff(4) filesize(4)
//	maxprot(4) initprot(4) nsects(4) flags(4); section sections, 68 bytes each
type Segment struct {
	cmd      Cmd
	cmdsize  uint32
	SegName  string
	Vmaddr   uint64
	Vmsize   uint64
	Fileoff  uint64
	Filesize uint64
	Maxprot  uint32
	Initprot uint32
	Nsects   uint32
	Flags    SegFlag

	sections []*Section
}

func NewSegment(cmd Cmd, cmdsize uint32) *Segment {
	return &Segment{
		cmd:     cmd,
		cmdsize: cmdsize,
	}
}

func (s *Segment) Cmd() Cmd {
	return s.cmd
}
func (s *Segment) Cmdsize() uint32 {
	return s.cmdsize
}

// Sections returns the sections of the segment.
func (s *Segment) Sections() []*Section {
	return s.sections
}

// Section is a full section_64 (80 bytes) or section (68 bytes).
//
// section_64:
//
//	sectname(16) segname(16) addr(8) size(8) offset(4) align(4) reloff(4)
//	nreloc(4) flags(4) reserved1(4) reserved2(4) reserved3(4)
//
// section: the same, but addr/size are 4 bytes each and there is no reserved3.
type Section struct {
	SectName  string
	SegName   string
	Addr      uint64
	Size      uint64
	Offset    uint32 // offset from the start of the Mach-O object (slice)
	Align     uint32 // log2 of the alignment
	Reloff    uint32 // offset of the section's relocation table
	Nreloc    uint32
	Flags     uint32 // low byte is SectionType, high bits are AttrFlag
	Reserved1 uint32 // e.g. an index into the indirect symtab
	Reserved2 uint32 // e.g. the stub size
	Reserved3 uint32

	segment *Segment
	f       *File
	data    []byte // Data() cache
}

func NewSection(segment *Segment, f *File) *Section {
	return &Section{
		segment: segment,
		f:       f,
	}
}

// Type returns the section type (the low byte of flags).
func (s *Section) Type() SectionType {
	return SectionType(s.Flags & 0xff)
}

// Attr returns the section attributes (the high bits of flags).
func (s *Section) Attr() AttrFlag {
	return AttrFlag(s.Flags & 0xffffff00)
}

// Segment returns the segment the section belongs to.
func (s *Section) Segment() *Segment {
	return s.segment
}

// Data returns the contents of the section: read lazily from the buffer
// and cached. Section offsets of a FAT slice are relative to the start of
// the slice, so we read from f.base. For ZEROFILL sections it returns zero
// bytes (there is no data in the file).
func (s *Section) Data() ([]byte, error) {
	if s.data != nil {
		return s.data, nil
	}

	t := s.Type()
	if t == S_ZEROFILL || t == S_GB_ZEROFILL || t == S_THREAD_LOCAL_ZEROFILL || s.Size == 0 {
		s.data = []byte{}
		return s.data, nil
	}

	data, err := readBytes(s.f.buf, s.f.base+int(s.Offset), int(s.Size))
	if err != nil {
		return nil, err
	}

	s.data = data
	return data, nil
}

// segment/section layouts.
const (
	segment64Size   = 72
	section64Size   = 80
	segment64Nsects = 64 // nsects offset within the command
	segment64Secs   = 72 // offset of the section_64 array

	segment32Size   = 56
	section32Size   = 68
	segment32Nsects = 48
	segment32Secs   = 56
)

// Segments returns all segments of the file (LC_SEGMENT_64 / LC_SEGMENT).
func (f *File) Segments() []*Segment {
	return f.segments
}

// Sections returns all sections of the file (flat across all segments).
func (f *File) Sections() []*Section {
	return f.sections
}

// Section returns a section by name (e.g. "__text") or nil.
func (f *File) Section(name string) *Section {
	for _, s := range f.sections {
		if s.SectName == name {
			return s
		}
	}

	return nil
}

// parseSegment parses a segment(_64) and its sections at absolute offset off.
func (f *File) parseSegment(off int, cmd Cmd, cmdsize uint32) (LoadCommand, error) {
	seg := NewSegment(cmd, cmdsize)
	is64 := cmd == LC_SEGMENT_64

	var err error
	nameRaw, err := readBytes(f.buf, off+8, 16)
	if err != nil {
		return nil, err
	}

	seg.SegName = cstr(nameRaw)
	if is64 {
		if seg.Vmaddr, err = readU64At(f.buf, f.order, off+24); err != nil {
			return nil, err
		}

		if seg.Vmsize, err = readU64At(f.buf, f.order, off+32); err != nil {
			return nil, err
		}

		if seg.Fileoff, err = readU64At(f.buf, f.order, off+40); err != nil {
			return nil, err
		}

		if seg.Filesize, err = readU64At(f.buf, f.order, off+48); err != nil {
			return nil, err
		}

		if seg.Maxprot, err = readU32At(f.buf, f.order, off+56); err != nil {
			return nil, err
		}

		if seg.Initprot, err = readU32At(f.buf, f.order, off+60); err != nil {
			return nil, err
		}

		if seg.Nsects, err = readU32At(f.buf, f.order, off+segment64Nsects); err != nil {
			return nil, err
		}

		var flagsRaw uint32
		if flagsRaw, err = readU32At(f.buf, f.order, off+68); err != nil {
			return nil, err
		}

		seg.Flags = SegFlag(flagsRaw)

		secBase := off + segment64Secs
		for i := range seg.Nsects {
			s := secBase + int(i*section64Size)
			sec, err := f.parseSection64(s, seg)
			if err != nil {
				return nil, err
			}

			seg.sections = append(seg.sections, sec)
			f.sections = append(f.sections, sec)
		}

		f.segments = append(f.segments, seg)
		return seg, nil
	}

	// 32-bit vmaddr/vmsize/fileoff/filesize are 4 bytes each.
	var v uint32
	if v, err = readU32At(f.buf, f.order, off+24); err != nil {
		return nil, err
	}

	seg.Vmaddr = uint64(v)
	if v, err = readU32At(f.buf, f.order, off+28); err != nil {
		return nil, err
	}

	seg.Vmsize = uint64(v)
	if v, err = readU32At(f.buf, f.order, off+32); err != nil {
		return nil, err
	}

	seg.Fileoff = uint64(v)
	if v, err = readU32At(f.buf, f.order, off+36); err != nil {
		return nil, err
	}

	seg.Filesize = uint64(v)
	if seg.Maxprot, err = readU32At(f.buf, f.order, off+40); err != nil {
		return nil, err
	}

	if seg.Initprot, err = readU32At(f.buf, f.order, off+44); err != nil {
		return nil, err
	}

	if seg.Nsects, err = readU32At(f.buf, f.order, off+segment32Nsects); err != nil {
		return nil, err
	}

	var flagsRaw uint32
	if flagsRaw, err = readU32At(f.buf, f.order, off+52); err != nil {
		return nil, err
	}

	seg.Flags = SegFlag(flagsRaw)

	secBase := off + segment32Secs
	for i := range seg.Nsects {
		s := secBase + int(i*section32Size)
		sec := NewSection(seg, f)

		nameRaw, err := readBytes(f.buf, s, 16)
		if err != nil {
			return nil, err
		}

		sec.SectName = cstr(nameRaw)
		if nameRaw, err = readBytes(f.buf, s+16, 16); err != nil {
			return nil, err
		}

		sec.SegName = cstr(nameRaw)
		if v, err = readU32At(f.buf, f.order, s+32); err != nil {
			return nil, err
		}

		sec.Addr = uint64(v)
		if v, err = readU32At(f.buf, f.order, s+36); err != nil {
			return nil, err
		}

		sec.Size = uint64(v)
		if sec.Offset, err = readU32At(f.buf, f.order, s+40); err != nil {
			return nil, err
		}

		if sec.Align, err = readU32At(f.buf, f.order, s+44); err != nil {
			return nil, err
		}

		if sec.Reloff, err = readU32At(f.buf, f.order, s+48); err != nil {
			return nil, err
		}

		if sec.Nreloc, err = readU32At(f.buf, f.order, s+52); err != nil {
			return nil, err
		}

		if sec.Flags, err = readU32At(f.buf, f.order, s+56); err != nil {
			return nil, err
		}

		if sec.Reserved1, err = readU32At(f.buf, f.order, s+60); err != nil {
			return nil, err
		}

		if sec.Reserved2, err = readU32At(f.buf, f.order, s+64); err != nil {
			return nil, err
		}

		seg.sections = append(seg.sections, sec)
		f.sections = append(f.sections, sec)
	}

	f.segments = append(f.segments, seg)
	return seg, nil
}

// parseSection64 parses a section_64 at absolute offset s.
func (f *File) parseSection64(s int, seg *Segment) (*Section, error) {
	sec := NewSection(seg, f)

	nameRaw, err := readBytes(f.buf, s, 16)
	if err != nil {
		return nil, err
	}

	sec.SectName = cstr(nameRaw)
	if nameRaw, err = readBytes(f.buf, s+16, 16); err != nil {
		return nil, err
	}

	sec.SegName = cstr(nameRaw)
	if sec.Addr, err = readU64At(f.buf, f.order, s+32); err != nil {
		return nil, err
	}

	if sec.Size, err = readU64At(f.buf, f.order, s+40); err != nil {
		return nil, err
	}

	if sec.Offset, err = readU32At(f.buf, f.order, s+48); err != nil {
		return nil, err
	}

	if sec.Align, err = readU32At(f.buf, f.order, s+52); err != nil {
		return nil, err
	}

	if sec.Reloff, err = readU32At(f.buf, f.order, s+56); err != nil {
		return nil, err
	}

	if sec.Nreloc, err = readU32At(f.buf, f.order, s+60); err != nil {
		return nil, err
	}

	if sec.Flags, err = readU32At(f.buf, f.order, s+64); err != nil {
		return nil, err
	}

	if sec.Reserved1, err = readU32At(f.buf, f.order, s+68); err != nil {
		return nil, err
	}

	if sec.Reserved2, err = readU32At(f.buf, f.order, s+72); err != nil {
		return nil, err
	}

	if sec.Reserved3, err = readU32At(f.buf, f.order, s+76); err != nil {
		return nil, err
	}

	return sec, nil
}
