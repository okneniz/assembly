package elf

// GNU symbol versioning (.gnu.version / .gnu.version_r / .gnu.version_d):
// links .dynsym entries with version names from external libraries.

// VersionNeedEntry is a single symbol that requires a version from a provider
// library.
type VersionNeedEntry struct {
	Idx   uint16 // version index (matches the .gnu.version entry)
	Flags uint16 // VER_FLG_WEAK and the like
	Name  string
}

func NewVersionNeedEntry(idx uint16, flags uint16, name string) VersionNeedEntry {
	return VersionNeedEntry{
		Idx:   idx,
		Flags: flags,
		Name:  name,
	}
}

// VersionNeed is a verneed entry: a library providing versions.
type VersionNeed struct {
	Version uint16 // version of the verneed structure (1)
	Cnt     uint16
	File    string
	Entries []VersionNeedEntry
}

func NewVersionNeed(version uint16, cnt uint16, file string) VersionNeed {
	return VersionNeed{
		Version: version,
		Cnt:     cnt,
		File:    file,
	}
}

// VersionDef is a verdef entry: a version defined by the object itself.
type VersionDef struct {
	Version uint16 // version of the verdef structure (1)
	Flags   uint16
	Idx     uint16
	Hash    uint32
	Name    string
	Parents []string
}

func NewVersionDef(version uint16, flags uint16, idx uint16, hash uint32) VersionDef {
	return VersionDef{
		Version: version,
		Flags:   flags,
		Idx:     idx,
		Hash:    hash,
	}
}

// Versioning flags.
const (
	VER_FLG_BASE uint16 = 0x1
	VER_FLG_WEAK uint16 = 0x2
	VER_FLG_INFO uint16 = 0x4
)

// VersionNeeds returns the .gnu.version_r table (which versions the object needs).
func (f *File) VersionNeeds() ([]VersionNeed, error) {
	data, err := f.versionData(SHT_GNU_VERNEED)
	if data == nil {
		return nil, err
	}

	// Elf64_Verneed = Elf32_Verneed (16 bytes):
	//	vn_version(2) vn_cnt(2) vn_file(4) vn_aux(4) vn_next(4)
	// Elf64_Vernaux = Elf32_Vernaux (16 bytes):
	//	vna_hash(4) vna_flags(2) vna_other(2) vna_name(4) vna_next(4)
	names, err := f.dynStrData()
	if err != nil {
		return nil, err
	}

	var out []VersionNeed
	for off := 0; off+16 <= len(data); {
		vnVersion := u16(data[off:], f.order)
		vnCnt := u16(data[off+2:], f.order)
		vnFile := u32(data[off+4:], f.order)
		vnAux := u32(data[off+8:], f.order)
		vnNext := u32(data[off+12:], f.order)

		need := NewVersionNeed(vnVersion, vnCnt, dynStr(names, vnFile))
		auxOff := off + int(vnAux)
		for range vnCnt {
			if auxOff+16 > len(data) {
				return nil, errf("verneed: vernaux is out of bounds")
			}

			vnaFlags := u16(data[auxOff+4:], f.order)
			vnaOther := u16(data[auxOff+6:], f.order)
			vnaName := u32(data[auxOff+8:], f.order)
			vnaNext := u32(data[auxOff+12:], f.order)

			need.Entries = append(
				need.Entries,
				NewVersionNeedEntry(vnaOther, vnaFlags, dynStr(names, vnaName)),
			)
			if vnaNext == 0 {
				break
			}

			auxOff += int(vnaNext)
		}

		out = append(out, need)
		if vnNext == 0 {
			break
		}

		off += int(vnNext)
	}

	return out, nil
}

// VersionDefs returns the .gnu.version_d table (versions defined by the object).
func (f *File) VersionDefs() ([]VersionDef, error) {
	data, err := f.versionData(SHT_GNU_VERDEF)
	if data == nil {
		return nil, err
	}

	// Elf64_Verdef = Elf32_Verdef (20 bytes):
	//	vd_version(2) vd_flags(2) vd_ndx(2) vd_cnt(2) vd_hash(4)
	//	vd_aux(4) vd_next(4)
	// Elf64_Verdaux (8 bytes): vda_name(4) vda_next(4)
	names, err := f.dynStrData()
	if err != nil {
		return nil, err
	}

	var out []VersionDef
	for off := 0; off+20 <= len(data); {
		vdVersion := u16(data[off:], f.order)
		vdFlags := u16(data[off+2:], f.order)
		vdNdx := u16(data[off+4:], f.order)
		vdCnt := u16(data[off+6:], f.order)
		vdHash := u32(data[off+8:], f.order)
		vdAux := u32(data[off+12:], f.order)
		vdNext := u32(data[off+16:], f.order)

		def := NewVersionDef(vdVersion, vdFlags, vdNdx, vdHash)
		// The first name is the version itself, the rest are its parents.
		auxOff := off + int(vdAux)
		for i := range vdCnt {
			if auxOff+8 > len(data) {
				return nil, errf("verdef: verdaux is out of bounds")
			}

			vdaName := u32(data[auxOff:], f.order)
			vdaNext := u32(data[auxOff+4:], f.order)
			if i == 0 {
				def.Name = dynStr(names, vdaName)
			} else {
				def.Parents = append(def.Parents, dynStr(names, vdaName))
			}

			if vdaNext == 0 {
				break
			}

			auxOff += int(vdaNext)
		}

		out = append(out, def)
		if vdNext == 0 {
			break
		}

		off += int(vdNext)
	}

	return out, nil
}

// SymbolVersions returns the .gnu.version array - one uint16 per .dynsym
// entry (0/1 = local/global without a version, >1 is a version index).
func (f *File) SymbolVersions() ([]uint16, error) {
	data, err := f.versionData(SHT_GNU_VERSYM)
	if data == nil {
		return nil, err
	}

	out := make([]uint16, 0, len(data)/2)
	for off := 0; off+2 <= len(data); off += 2 {
		out = append(out, u16(data[off:], f.order))
	}

	return out, nil
}

// versionData reads the data of a versioning section; (nil, nil) means there
// is no section.
func (f *File) versionData(typ SectionType) ([]byte, error) {
	for _, s := range f.sections {
		if s.Type == typ && s.Size > 0 {
			return s.Data()
		}
	}

	return nil, nil
}

// dynStrData reads the dynamic string table (DT_STRTAB) whole.
func (f *File) dynStrData() ([]byte, error) {
	base, ok := f.dynamicValue(DT_STRTAB)
	if !ok {
		return nil, errf("no DT_STRTAB for version tables")
	}

	size, _ := f.dynamicValue(DT_STRSZ)
	if size == 0 {
		size = 0x10000 // a safeguard for files without DT_STRSZ
	}

	off, err := f.vaddrToOff(base)
	if err != nil {
		return nil, err
	}

	return readBytes(f.buf, int(off), int(size))
}

// dynStr extracts the string at an offset from the loaded strtab.
func dynStr(tab []byte, off uint32) string {
	if int(off) >= len(tab) {
		return ""
	}

	return cstr(tab[off:])
}
