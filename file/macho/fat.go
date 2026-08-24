package macho

import (
	"encoding/binary"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
)

// FatArch is a single slice of a FAT/Universal archive: architecture
// metadata and a lazily opened object.
type FatArch struct {
	CpuType    int32
	CpuSubtype int32
	Offset     uint64 // offset of the slice in the file
	Size       uint64
	Align      uint32 // log2

	buf  common.Buffer[byte, int] // buffer of the archive
	file *File
}

func NewFatArch(buf common.Buffer[byte, int]) *FatArch {
	return &FatArch{buf: buf}
}

// Open returns the parsed Mach-O object of the slice (cached).
func (a *FatArch) Open() (*File, error) {
	if a.file != nil {
		return a.file, nil
	}

	f, err := openAt(a.buf, int(a.Offset))
	if err != nil {
		return nil, err
	}

	a.file = f
	return f, nil
}

// Fat is a whole FAT/Universal archive: all slices.
type Fat struct {
	buf    common.Buffer[byte, int]
	Arches []*FatArch
}

func NewFat(buf common.Buffer[byte, int]) *Fat {
	return &Fat{buf: buf}
}

// OpenFat opens a FAT/Universal archive with all slices. For a plain
// Mach-O it returns an error.
func OpenFat(path string) (*Fat, error) {
	buf, err := bytes.BufferFromFile(path)
	if err != nil {
		return nil, err
	}

	magicLE, err := readU32At(buf, binary.LittleEndian, 0)
	if err != nil {
		return nil, err
	}

	if !isFAT(magicLE) {
		return nil, errf("file is not a FAT/Universal archive (magic LE %#x)", magicLE)
	}

	var order binary.ByteOrder = binary.LittleEndian
	switch magicLE {
	case FAT_CIGAM, FAT_CIGAM_64:
		order = binary.BigEndian
	}

	is64 := magicLE == FAT_MAGIC_64 || magicLE == FAT_CIGAM_64
	entrySize := uint32(20)
	if is64 {
		entrySize = 32
	}

	nfat, err := readU32At(buf, order, 4)
	if err != nil {
		return nil, err
	}

	fat := NewFat(buf)
	for i := range nfat {
		base := 8 + int(i*entrySize)
		a := NewFatArch(buf)
		var raw uint32
		if raw, err = readU32At(buf, order, base); err != nil {
			return nil, err
		}

		a.CpuType = int32(raw)
		if raw, err = readU32At(buf, order, base+4); err != nil {
			return nil, err
		}

		a.CpuSubtype = int32(raw)
		if is64 {
			if a.Offset, err = readU64At(buf, order, base+8); err != nil {
				return nil, err
			}

			if a.Size, err = readU64At(buf, order, base+16); err != nil {
				return nil, err
			}

			if raw, err = readU32At(buf, order, base+24); err != nil {
				return nil, err
			}

			a.Align = raw
		} else {
			if raw, err = readU32At(buf, order, base+8); err != nil {
				return nil, err
			}

			a.Offset = uint64(raw)
			if raw, err = readU32At(buf, order, base+12); err != nil {
				return nil, err
			}

			a.Size = uint64(raw)
			if raw, err = readU32At(buf, order, base+16); err != nil {
				return nil, err
			}

			a.Align = raw
		}

		fat.Arches = append(fat.Arches, a)
	}

	return fat, nil
}

// Arch returns the slice with the given cputype, or nil.
func (f *Fat) Arch(cpu int32) *FatArch {
	for _, a := range f.Arches {
		if a.CpuType == cpu {
			return a
		}
	}

	return nil
}
