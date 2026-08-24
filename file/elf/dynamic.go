package elf

// DynamicEntry is an entry of the .dynamic section / PT_DYNAMIC segment:
// d_tag (which parameter it is) + d_val/d_ptr (the value; for pointers, a
// virtual address that must be mapped to a file offset via PT_LOAD).
//
// Elf64_Dyn (16 bytes): d_tag(8) d_val(8). Elf32_Dyn (8 bytes): d_tag(4) d_val(4).
type DynamicEntry struct {
	Tag DynTag
	Val uint64
}

func NewDynamicEntry(tag DynTag, val uint64) DynamicEntry {
	return DynamicEntry{
		Tag: tag,
		Val: val,
	}
}

const (
	dyn64Size = 16
	dyn32Size = 8
)

// Dynamic returns the .dynamic entries: the PT_DYNAMIC segment, or, in its
// absence (ET_REL), the SHT_DYNAMIC section. nil without an error means a file
// without dynamics. Entries are read up to DT_NULL. The result is cached.
func (f *File) Dynamic() ([]DynamicEntry, error) {
	if f.dynDone {
		return f.dyn, f.dynErr
	}

	f.dynDone = true

	data, err := f.dynamicData()
	if err != nil {
		f.dynErr = err
		return nil, err
	}

	if data == nil {
		return nil, nil // no dynamics - a normal situation (ET_REL/static)
	}

	size := dyn32Size
	if f.class == CLASS64 {
		size = dyn64Size
	}

	var out []DynamicEntry
	for off := 0; off+size <= len(data); off += size {
		var tag, val uint64
		if f.class == CLASS64 {
			tag = u64(data[off:], f.order)
			val = u64(data[off+8:], f.order)
		} else {
			tag = uint64(u32(data[off:], f.order))
			val = uint64(u32(data[off+4:], f.order))
		}

		out = append(out, NewDynamicEntry(DynTag(tag), val))
		if tag == uint64(DT_NULL) {
			break
		}
	}

	f.dyn = out
	return out, nil
}

// dynamicData reads the raw dynamics bytes from PT_DYNAMIC or SHT_DYNAMIC.
func (f *File) dynamicData() ([]byte, error) {
	for _, p := range f.progs {
		if p.Type != PT_DYNAMIC || p.Filesz == 0 {
			continue
		}

		return readBytes(f.buf, int(p.Off), int(p.Filesz))
	}

	for _, s := range f.sections {
		if s.Type == SHT_DYNAMIC && s.Size > 0 {
			return s.Data()
		}
	}

	return nil, nil
}

// dynamicValue returns the value of the first tag entry (or 0/false).
func (f *File) dynamicValue(tag DynTag) (uint64, bool) {
	entries, err := f.Dynamic()
	if err != nil {
		return 0, false
	}

	for _, e := range entries {
		if e.Tag == tag {
			return e.Val, true
		}
	}

	return 0, false
}

// Needed returns the names of shared libraries (DT_NEEDED).
func (f *File) Needed() ([]string, error) {
	entries, err := f.Dynamic()
	if err != nil {
		return nil, err
	}

	var out []string
	for _, e := range entries {
		if e.Tag == DT_NEEDED {
			s, err := f.dynStr(e.Val)
			if err != nil {
				return nil, err
			}

			out = append(out, s)
		}
	}

	return out, nil
}

// Soname returns DT_SONAME ("" if the tag is absent).
func (f *File) Soname() (string, error) {
	return f.dynTagStr(DT_SONAME)
}

// Rpath returns DT_RPATH ("" if the tag is absent; a deprecated tag).
func (f *File) Rpath() (string, error) {
	return f.dynTagStr(DT_RPATH)
}

// Runpath returns DT_RUNPATH ("" if the tag is absent).
func (f *File) Runpath() (string, error) {
	return f.dynTagStr(DT_RUNPATH)
}

// dynTagStr reads the string at a string tag's value (an offset into DT_STRTAB).
func (f *File) dynTagStr(tag DynTag) (string, error) {
	v, ok := f.dynamicValue(tag)
	if !ok {
		return "", nil
	}

	return f.dynStr(v)
}

// dynStr reads a null-terminated string from DT_STRTAB at offset off.
func (f *File) dynStr(off uint64) (string, error) {
	base, ok := f.dynamicValue(DT_STRTAB)
	if !ok {
		return "", errf("no DT_STRTAB for dynamic strings")
	}

	size, _ := f.dynamicValue(DT_STRSZ)
	fileOff, err := f.vaddrToOff(base)
	if err != nil {
		return "", err
	}

	if off >= size {
		return "", errf("string offset %#x is outside DT_STRTAB (size %#x)", off, size)
	}

	return readCStringAt(f.buf, int(fileOff+off))
}

// vaddrToOff maps a virtual address to a file offset via PT_LOAD.
func (f *File) vaddrToOff(vaddr uint64) (uint64, error) {
	for _, p := range f.progs {
		if p.Type != PT_LOAD {
			continue
		}

		if vaddr >= p.Vaddr && vaddr < p.Vaddr+p.Filesz {
			return p.Off + (vaddr - p.Vaddr), nil
		}
	}

	return 0, errf("vaddr %#x is not covered by any PT_LOAD", vaddr)
}
