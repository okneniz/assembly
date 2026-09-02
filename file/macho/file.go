package macho

import (
	"encoding/binary"

	"github.com/okneniz/parsec"
	"github.com/okneniz/parsec/bytes"
)

// Header is a full mach_header{,_64} with normalized fields.
//
// mach_header_64 (32 bytes): magic(4) cputype(4) cpusubtype(4) filetype(4)
//
//	ncmds(4) sizeofcmds(4) flags(4) reserved(4)
//
// mach_header (28 bytes): the same without reserved.
type Header struct {
	Magic      uint32
	CpuType    int32
	CpuSubtype int32
	FileType   FileType
	Ncmds      uint32
	Sizeofcmds uint32
	Flags      HeaderFlag
	Reserved   uint32 // 64-bit only
}

// Header size of the selected class (28 or 32).
const (
	headerSize64 = 32
	headerSize32 = 28
)

// File is a parsed Mach-O object. base is the object's offset in the file
// (for a FAT archive slice; 0 for a plain file). The header, load commands,
// and sections are read eagerly; the symbol tables, relocations, and dyld
// streams lazily.
type File struct {
	buf      parsec.Buffer[byte, int]
	base     int
	is64     bool
	order    binary.ByteOrder
	hdr      Header
	commands []LoadCommand
	segments []*Segment
	sections []*Section

	symtab     []Symbol
	symtabErr  error
	symtabDone bool
}

func NewFile(buf parsec.Buffer[byte, int], base int, is64 bool, order binary.ByteOrder) *File {
	return &File{
		buf:   buf,
		base:  base,
		is64:  is64,
		order: order,
	}
}

// IsMagic checks whether the first 4 bytes look like Mach-O (including FAT).
func IsMagic(m [4]byte) bool {
	le := binary.LittleEndian.Uint32(m[:])
	be := binary.BigEndian.Uint32(m[:])
	switch le {
	case MH_MAGIC,
		MH_CIGAM,
		MH_MAGIC_64,
		MH_CIGAM_64,
		FAT_MAGIC,
		FAT_CIGAM,
		FAT_MAGIC_64,
		FAT_CIGAM_64:
		return true
	}

	switch be {
	case MH_MAGIC,
		MH_CIGAM,
		MH_MAGIC_64,
		MH_CIGAM_64,
		FAT_MAGIC,
		FAT_CIGAM,
		FAT_MAGIC_64,
		FAT_CIGAM_64:
		return true
	}

	return false
}

// Open opens a Mach-O file. For a FAT/Universal archive it selects the
// arm64 slice (an error if there is none); for a plain Mach-O it parses it
// directly.
func Open(path string) (*File, error) {
	buf, err := bytes.BufferFromFile(path)
	if err != nil {
		return nil, err
	}

	magic, err := readU32At(buf, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}

	if isFAT(magic) {
		base, err := findFATSlice(buf, magic, CPU_TYPE_ARM64)
		if err != nil {
			return nil, err
		}

		return openAt(buf, base)
	}

	return openAt(buf, 0)
}

// openAt parses a Mach-O object located at offset base in the buffer.
func openAt(buf parsec.Buffer[byte, int], base int) (*File, error) {
	magic, err := readU32At(buf, binary.LittleEndian, base)
	if err != nil {
		return nil, err
	}

	var order binary.ByteOrder = binary.LittleEndian
	is64 := true
	switch magic {
	case MH_MAGIC_64:
	case MH_CIGAM_64:
		order = binary.BigEndian
	case MH_MAGIC:
		is64 = false
	case MH_CIGAM:
		is64 = false
		order = binary.BigEndian
	default:
		return nil, errf("no Mach-O magic at offset %d (read %#x)", base, magic)
	}

	f := NewFile(buf, base, is64, order)

	if err := f.parseHeader(); err != nil {
		return nil, err
	}

	if err := f.parseLoadCommands(); err != nil {
		return nil, err
	}

	return f, nil
}

// Header returns the full header of the object.
func (f *File) Header() Header {
	return f.hdr
}

// parseHeader reads the mach_header{,_64} at offset base.
func (f *File) parseHeader() error {
	var h Header
	// magic is read as LE: a cigam value means a file with reversed byte order.
	var err error
	if h.Magic, err = readU32At(f.buf, binary.LittleEndian, f.base); err != nil {
		return err
	}

	var raw uint32
	if raw, err = readU32At(f.buf, f.order, f.base+4); err != nil {
		return err
	}

	h.CpuType = int32(raw)
	if raw, err = readU32At(f.buf, f.order, f.base+8); err != nil {
		return err
	}

	h.CpuSubtype = int32(raw)
	if raw, err = readU32At(f.buf, f.order, f.base+12); err != nil {
		return err
	}

	h.FileType = FileType(raw)
	if h.Ncmds, err = readU32At(f.buf, f.order, f.base+16); err != nil {
		return err
	}

	if h.Sizeofcmds, err = readU32At(f.buf, f.order, f.base+20); err != nil {
		return err
	}

	if raw, err = readU32At(f.buf, f.order, f.base+24); err != nil {
		return err
	}

	h.Flags = HeaderFlag(raw)
	if f.is64 {
		if h.Reserved, err = readU32At(f.buf, f.order, f.base+28); err != nil {
			return err
		}
	}

	f.hdr = h
	return nil
}

// isFAT determines from magic (read as LE) whether the file is a FAT archive.
func isFAT(magic uint32) bool {
	switch magic {
	case FAT_MAGIC, FAT_CIGAM, FAT_MAGIC_64, FAT_CIGAM_64:
		return true
	}

	return false
}

// findFATSlice searches a FAT archive for the slice with the given cputype
// and returns its offset in the file. FAT header: magic(4) nfat_arch(4);
// fat_arch is 20 bytes (fat_arch_64 is 32): cputype cpusubtype offset size
// align.
func findFATSlice(buf parsec.Buffer[byte, int], magicLE uint32, cpu int32) (int, error) {
	// magicLE is the magic read as LE: a standard BE archive (CA FE BA BE)
	// yields FAT_CIGAM -> the header fields are read as BigEndian, and vice
	// versa.
	var order binary.ByteOrder = binary.LittleEndian
	switch magicLE {
	case FAT_CIGAM, FAT_CIGAM_64:
		order = binary.BigEndian
	}

	nfat, err := readU32At(buf, order, 4)
	if err != nil {
		return 0, err
	}

	is64 := magicLE == FAT_MAGIC_64 || magicLE == FAT_CIGAM_64
	entrySize := uint32(20)
	if is64 {
		entrySize = 32
	}

	for i := range nfat {
		base := 8 + int(i*entrySize)
		cputype32, err := readU32At(buf, order, base)
		if err != nil {
			return 0, err
		}

		if int32(cputype32) != cpu {
			continue
		}

		var off uint64
		if is64 {
			if off, err = readU64At(buf, order, base+8); err != nil {
				return 0, err
			}
		} else {
			o32, err2 := readU32At(buf, order, base+8)
			if err2 != nil {
				return 0, err2
			}

			off = uint64(o32)
		}

		return int(off), nil
	}

	return 0, errf("no %s architecture in FAT archive", cpuName(cpu))
}
