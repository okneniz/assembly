package elf

// A minimal static ELF64 emitter: wraps sections into an executable file that
// can be run natively on Linux arm64/riscv64 or under qemu-user
// (qemu-aarch64 / qemu-riscv64) on any host - including macOS.
// Writing an ELF lives here (next to the parser): the assembly assembler knows
// nothing about formats.
//
// File layout:
//
//	0x0000   Elf64_Ehdr (64) + zero padding (ehdr does not have to be mapped)
//	0x10000  the section image (as is, back to back; with an unaligned base -
//	         prefixed with zeros up to base)
//	at the end  Elf64_Phdr - INSIDE the mapped area of the segment
//
// The single PT_LOAD maps the image at a 64 KiB-aligned vaddr <= base (the
// code lands exactly at base - absolute la/call addressing stays correct),
// and the phdr sits at the end of the image with e_phoff pointing at it: the
// kernel builds auxv AT_PHDR as vaddr+e_phoff, and the runtimes of real ELF
// programs (e.g. Go) read the program headers at that address.
//
// 64 KiB (not 4 KiB) is the unit of BOTH the image offset and p_align: the
// kernel requires p_offset to be congruent to p_vaddr modulo the TARGET page
// size, and the supported targets run kernels with 4/16/64 KiB pages (the
// Alpine loongarch64 lts kernel, for one, is 16 KiB - a 0x1000 offset is
// rejected there with EINVAL at execve).

import (
	"encoding/binary"
	"fmt"
)

const (
	elfTypeExec = 2 // ET_EXEC
	elfVersion  = 1

	ptLoad     = 1 // PT_LOAD
	pfRXW      = 7 // PF_R|PF_W|PF_X - demo loader: code+data in one page
	elfCodeOff = 0x10000
)

// Blob is a section to write: data Data at virtual address Addr
// (names are not preserved in the minimal format). MemSize is the size in
// memory: the difference MemSize-len(Data) is a zero NOBITS tail (.bss) that
// occupies no file space; it is allowed only on the last blob (otherwise the
// file layout would diverge from the vaddr layout).
type Blob struct {
	Addr    uint64
	Data    []byte
	MemSize int
}

func NewBlob(addr uint64, data []byte) Blob {
	return Blob{
		Addr:    addr,
		Data:    data,
		MemSize: len(data),
	}
}

// NewNobitsBlob is a blob with a zero tail in memory: the file part is data,
// the total size is memSize.
func NewNobitsBlob(addr uint64, data []byte, memSize int) Blob {
	return Blob{
		Addr:    addr,
		Data:    data,
		MemSize: memSize,
	}
}

// Write assembles an executable ELF64 LE from sections. machine is e_machine
// (EM_AARCH64/EM_RISCV); entry is the absolute address of the entry point
// (usually base + the offset of the start symbol). Sections are placed into
// the file back to back.
//
// The base is rounded down to a 64 KiB boundary (p_vaddr must be aligned
// and congruent to p_offset for every target page size; the image is padded
// with zeros up to base), and the
// phdr is placed at the end of the image - inside the mapped area of the
// segment (see the package comment about AT_PHDR). The NOBITS tail of the last
// blob yields p_memsz > p_filesz: the file carries no zeros, the kernel
// reserves them as anonymous zero pages.
func Write(machine Machine, flags uint32, base, entry uint64, blobs []Blob) ([]byte, error) {
	var code []byte
	memTail := 0
	for i, b := range blobs {
		if b.MemSize < len(b.Data) {
			return nil, fmt.Errorf("blob %d: MemSize %d < file size %d", i, b.MemSize, len(b.Data))
		}

		if tail := b.MemSize - len(b.Data); tail > 0 && i != len(blobs)-1 {
			return nil, fmt.Errorf("blob %d: NOBITS tail is only permitted on the last blob", i)
		}

		code = append(code, b.Data...)
		if i == len(blobs)-1 {
			memTail = b.MemSize - len(b.Data)
		}
	}

	imgBase := base &^ 0xffff
	pad := int(base - imgBase)
	image := make([]byte, pad+len(code))
	copy(image[pad:], code)

	phOff := elfCodeOff + len(image) // phdr right after the image
	out := make([]byte, phOff+56)

	// --- Elf64_Ehdr ---
	ehdr := out[:64]
	ehdr[0] = Mag0
	copy(ehdr[1:4], "ELF")
	ehdr[EI_CLASS] = byte(CLASS64)
	ehdr[EI_DATA] = byte(DATA2LSB)
	ehdr[EI_VERSION] = byte(elfVersion)
	binary.LittleEndian.PutUint16(ehdr[hdr64Type:], elfTypeExec)
	binary.LittleEndian.PutUint16(ehdr[hdr64Machine:], uint16(machine))
	binary.LittleEndian.PutUint32(ehdr[hdr64Flags:], flags)
	binary.LittleEndian.PutUint32(ehdr[hdr64Version:], elfVersion)
	binary.LittleEndian.PutUint64(ehdr[hdr64Entry:], entry)
	binary.LittleEndian.PutUint64(ehdr[hdr64Phoff:], uint64(phOff))
	binary.LittleEndian.PutUint16(ehdr[hdr64Ehsize:], 64)
	binary.LittleEndian.PutUint16(ehdr[hdr64Phentsz:], 56)
	binary.LittleEndian.PutUint16(ehdr[hdr64Phnum:], 1)

	// --- Elf64_Phdr (a single PT_LOAD; offset 0x10000 -> vaddr = imgBase) ---
	phdr := out[phOff:]
	binary.LittleEndian.PutUint32(phdr[0:], ptLoad)
	binary.LittleEndian.PutUint32(phdr[4:], pfRXW)
	binary.LittleEndian.PutUint64(phdr[8:], elfCodeOff)                     // p_offset
	binary.LittleEndian.PutUint64(phdr[16:], imgBase)                       // p_vaddr
	binary.LittleEndian.PutUint64(phdr[24:], imgBase)                       // p_paddr
	binary.LittleEndian.PutUint64(phdr[32:], uint64(len(image)+56))         // p_filesz
	binary.LittleEndian.PutUint64(phdr[40:], uint64(len(image)+56+memTail)) // p_memsz
	binary.LittleEndian.PutUint64(phdr[48:], 0x10000)                       // p_align

	copy(out[elfCodeOff:], image)
	return out, nil
}
