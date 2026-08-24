package elf

import (
	"encoding/binary"

	"github.com/okneniz/parsec/bytes"
	"github.com/okneniz/parsec/common"
)

// Header is a full Elf{32,64}_Ehdr with normalized fields: all offsets and
// sizes are widened to uint64, and the byte order is a ready binary.ByteOrder.
// The Phnum/Shnum/Shstrndx counters are already resolved according to the
// extended numbering rules (the real values from shdr[0] when the 16-bit
// fields overflow).
type Header struct {
	Ident      [EI_NIDENT]byte
	Class      Class
	Order      binary.ByteOrder
	OSABI      OSABI
	ABIVersion uint8
	Type       Type
	Machine    Machine
	Version    uint32 // e_version (EV_CURRENT = 1)
	Entry      uint64
	Phoff      uint64
	Shoff      uint64
	Flags      uint32
	Ehsize     uint16
	Phentsize  uint16
	Shentsize  uint16
	Phnum      int
	Shnum      int
	Shstrndx   int
}

// Field offsets of Elf64_Ehdr.
const (
	hdr64Type    = 16
	hdr64Machine = 18
	hdr64Version = 20
	hdr64Entry   = 24
	hdr64Phoff   = 32
	hdr64Shoff   = 40
	hdr64Flags   = 48
	hdr64Ehsize  = 52
	hdr64Phentsz = 54
	hdr64Phnum   = 56
	hdr64Shentsz = 58
	hdr64Shnum   = 60
	hdr64Shstrnx = 62
)

// Field offsets of Elf32_Ehdr.
const (
	hdr32Type    = 16
	hdr32Machine = 18
	hdr32Version = 20
	hdr32Entry   = 24
	hdr32Phoff   = 28
	hdr32Shoff   = 32
	hdr32Flags   = 36
	hdr32Ehsize  = 40
	hdr32Phentsz = 42
	hdr32Phnum   = 44
	hdr32Shentsz = 46
	hdr32Shnum   = 48
	hdr32Shstrnx = 50
)

// File is a parsed ELF: a single positional buffer for the whole file; the
// header, program headers, and section table are read eagerly; symbol,
// relocation, and other tables lazily via methods.
type File struct {
	buf      common.Buffer[byte, int]
	order    binary.ByteOrder
	class    Class
	hdr      Header
	progs    []Prog
	sections []*Section

	symtab     []Symbol
	symtabErr  error
	symtabDone bool

	dynsym     []Symbol
	dynsymErr  error
	dynsymDone bool

	dyn     []DynamicEntry
	dynErr  error
	dynDone bool
}

func NewFile(buf common.Buffer[byte, int], order binary.ByteOrder, class Class) *File {
	return &File{
		buf:   buf,
		order: order,
		class: class,
	}
}

// IsMagic checks the first 4 bytes of a file for the ELF magic number.
func IsMagic(m [4]byte) bool {
	return m[0] == Mag0 && m[1] == Mag1 && m[2] == Mag2 && m[3] == Mag3
}

// Open opens a file, validates the magic number/class/byte order, and parses
// the header, program headers, and section headers (including extended
// numbering). An error is returned for any invalid data; there are no panics.
func Open(path string) (*File, error) {
	buf, err := bytes.BufferFromFile(path)
	if err != nil {
		return nil, err
	}

	var magic [4]byte
	for i := range 4 {
		b, err := readByteAt(buf, EI_MAG0+i)
		if err != nil {
			return nil, err
		}

		magic[i] = b
	}

	if !IsMagic(magic) {
		return nil, errf("not an ELF file (magic %x)", magic)
	}

	// e_ident[5] = EI_DATA -> byte order; e_ident[4] = EI_CLASS -> 32/64.
	dataByte, err := readByteAt(buf, EI_DATA)
	if err != nil {
		return nil, err
	}

	var order binary.ByteOrder
	switch dataByte {
	case DATA2LSB:
		order = binary.LittleEndian
	case DATA2MSB:
		order = binary.BigEndian
	default:
		return nil, errf("unknown EI_DATA=%d", dataByte)
	}

	classByte, err := readByteAt(buf, EI_CLASS)
	if err != nil {
		return nil, err
	}

	class := Class(classByte)
	if class != CLASS32 && class != CLASS64 {
		return nil, errf("unknown EI_CLASS=%d", classByte)
	}

	f := NewFile(buf, order, class)

	if err := f.parseHeader(); err != nil {
		return nil, err
	}

	if err := f.parseProgs(); err != nil {
		return nil, err
	}

	if err := f.parseSections(); err != nil {
		return nil, err
	}

	return f, nil
}

// Header returns the full file header.
func (f *File) Header() Header {
	return f.hdr
}

// parseHeader reads all Elf{32,64}_Ehdr fields; counters that have escaped
// into reserved values (extended numbering) are resolved via shdr[0]:
//
//	e_shnum == 0 && e_shoff != 0 -> the real section count is in shdr[0].sh_size
//	e_shstrndx == SHN_XINDEX     -> the real index is in shdr[0].sh_link
//	e_phnum == PN_XNUM           -> the real phdr count is in shdr[0].sh_info
func (f *File) parseHeader() error {
	is64 := f.class == CLASS64
	var h Header
	h.Class = f.class
	h.Order = f.order

	for i := range EI_NIDENT {
		b, err := readByteAt(f.buf, i)
		if err != nil {
			return err
		}

		h.Ident[i] = b
	}

	h.OSABI = OSABI(h.Ident[EI_OSABI])
	h.ABIVersion = h.Ident[EI_ABIVERSION]

	var typeRaw, machineRaw uint16
	var versionRaw, flagsRaw uint32
	var phentsizeRaw, shentsizeRaw uint16
	var phnumRaw, shnumRaw, shstrndxRaw uint16
	if is64 {
		var err error
		if typeRaw, err = readU16At(f.buf, f.order, hdr64Type); err != nil {
			return err
		}

		if machineRaw, err = readU16At(f.buf, f.order, hdr64Machine); err != nil {
			return err
		}

		if versionRaw, err = readU32At(f.buf, f.order, hdr64Version); err != nil {
			return err
		}

		if h.Entry, err = readU64At(f.buf, f.order, hdr64Entry); err != nil {
			return err
		}

		if h.Phoff, err = readU64At(f.buf, f.order, hdr64Phoff); err != nil {
			return err
		}

		if h.Shoff, err = readU64At(f.buf, f.order, hdr64Shoff); err != nil {
			return err
		}

		if flagsRaw, err = readU32At(f.buf, f.order, hdr64Flags); err != nil {
			return err
		}

		if h.Ehsize, err = readU16At(f.buf, f.order, hdr64Ehsize); err != nil {
			return err
		}

		if phentsizeRaw, err = readU16At(f.buf, f.order, hdr64Phentsz); err != nil {
			return err
		}

		if phnumRaw, err = readU16At(f.buf, f.order, hdr64Phnum); err != nil {
			return err
		}

		if shentsizeRaw, err = readU16At(f.buf, f.order, hdr64Shentsz); err != nil {
			return err
		}

		if shnumRaw, err = readU16At(f.buf, f.order, hdr64Shnum); err != nil {
			return err
		}

		if shstrndxRaw, err = readU16At(f.buf, f.order, hdr64Shstrnx); err != nil {
			return err
		}
	} else {
		var err error
		if typeRaw, err = readU16At(f.buf, f.order, hdr32Type); err != nil {
			return err
		}

		if machineRaw, err = readU16At(f.buf, f.order, hdr32Machine); err != nil {
			return err
		}

		if versionRaw, err = readU32At(f.buf, f.order, hdr32Version); err != nil {
			return err
		}

		var entry, phoff, shoff uint32
		if entry, err = readU32At(f.buf, f.order, hdr32Entry); err != nil {
			return err
		}

		if phoff, err = readU32At(f.buf, f.order, hdr32Phoff); err != nil {
			return err
		}

		if shoff, err = readU32At(f.buf, f.order, hdr32Shoff); err != nil {
			return err
		}

		h.Entry, h.Phoff, h.Shoff = uint64(entry), uint64(phoff), uint64(shoff)
		if flagsRaw, err = readU32At(f.buf, f.order, hdr32Flags); err != nil {
			return err
		}

		if h.Ehsize, err = readU16At(f.buf, f.order, hdr32Ehsize); err != nil {
			return err
		}

		if phentsizeRaw, err = readU16At(f.buf, f.order, hdr32Phentsz); err != nil {
			return err
		}

		if phnumRaw, err = readU16At(f.buf, f.order, hdr32Phnum); err != nil {
			return err
		}

		if shentsizeRaw, err = readU16At(f.buf, f.order, hdr32Shentsz); err != nil {
			return err
		}

		if shnumRaw, err = readU16At(f.buf, f.order, hdr32Shnum); err != nil {
			return err
		}

		if shstrndxRaw, err = readU16At(f.buf, f.order, hdr32Shstrnx); err != nil {
			return err
		}
	}

	h.Type = Type(typeRaw)
	h.Machine = Machine(machineRaw)
	h.Version = versionRaw
	h.Flags = flagsRaw
	h.Phentsize = phentsizeRaw
	h.Shentsize = shentsizeRaw

	const pnXnum = 0xffff
	h.Phnum = int(phnumRaw)
	h.Shnum = int(shnumRaw)
	h.Shstrndx = int(shstrndxRaw)

	// Extended numbering: the real values live in shdr[0]. f.hdr is filled in
	// BEFORE parsing shdr[0] - readShdr0 takes e_shoff from it.
	f.hdr = h
	needShdr0 := (shnumRaw == 0 && h.Shoff != 0) ||
		shstrndxRaw == uint16(SHN_XINDEX) ||
		phnumRaw == pnXnum
	if needShdr0 {
		size, link, info, err := f.readShdr0()
		if err != nil {
			return err
		}

		if shnumRaw == 0 && h.Shoff != 0 {
			f.hdr.Shnum = int(size)
		}

		if shstrndxRaw == uint16(SHN_XINDEX) {
			f.hdr.Shstrndx = int(link)
		}

		if phnumRaw == pnXnum {
			f.hdr.Phnum = int(info)
		}
	}

	return nil
}

// readShdr0 reads sh_size/sh_link/sh_info of the zeroth section header (for
// extended numbering).
func (f *File) readShdr0() (size uint64, link, info uint32, retErr error) {
	base := int(f.hdr.Shoff)
	if f.class == CLASS64 {
		if size, retErr = readU64At(f.buf, f.order, base+32); retErr != nil {
			return
		}

		if link, retErr = readU32At(f.buf, f.order, base+40); retErr != nil {
			return
		}

		info, retErr = readU32At(f.buf, f.order, base+44)
		return
	}

	var s32, l32, i32 uint32
	if s32, retErr = readU32At(f.buf, f.order, base+20); retErr != nil {
		return
	}

	if l32, retErr = readU32At(f.buf, f.order, base+24); retErr != nil {
		return
	}

	if i32, retErr = readU32At(f.buf, f.order, base+28); retErr != nil {
		return
	}

	return uint64(s32), l32, i32, nil
}
