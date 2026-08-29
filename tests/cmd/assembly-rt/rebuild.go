package main

// Rebuilding the whole binary with our minimal ELF writer: the file-backed
// (PROGBITS) ALLOC sections of the original (read by our own parser) are laid
// out into one flat image byte-for-byte, .text is replaced with the
// round-tripped bytes; the NOBITS tail (.bss after the last file-backed
// section) is not materialized into the file - the writer emits
// p_memsz > p_filesz, and the kernel maps in the zeros. Zero external tools:
// the parser is ours, the assembler is ours, the emitter is ours.

import (
	"errors"
	"fmt"
	"os"

	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/file/elf"
)

// rebuildELF builds outPath from the original binPath: the memory image
// [minAddr, maxEnd) from the ALLOC sections, .text from the round-trip bytes
// text, entry is the original's entry point.
func rebuildELF(binPath string, text []byte, outPath string) error {
	ef, err := elf.Open(binPath)
	if err != nil {
		return err
	}

	hdr := ef.Header()

	machine, supported := elfMachines[hdr.Machine]
	if !supported {
		return fmt.Errorf("unsupported e_machine %v", hdr.Machine)
	}

	var minAddr, fileEnd, maxEnd uint64
	for _, s := range ef.Sections() {
		if !allocSection(s) {
			continue
		}

		if minAddr == 0 || s.Addr < minAddr {
			minAddr = s.Addr
		}

		if end := s.Addr + s.Size; end > maxEnd {
			maxEnd = end
		}

		if s.Type != elf.SHT_NOBITS {
			if end := s.Addr + s.Size; end > fileEnd {
				fileEnd = end
			}
		}
	}

	if minAddr == 0 || fileEnd <= minAddr {
		return errors.New("no file-backed ALLOC sections with addresses")
	}

	image := make([]byte, fileEnd-minAddr)
	for _, s := range ef.Sections() {
		if !allocSection(s) || s.Type == elf.SHT_NOBITS {
			continue
		}

		d, err := s.Data()
		if err != nil {
			return fmt.Errorf("section %s: %w", s.Name, err)
		}

		copy(image[s.Addr-minAddr:], d)
	}

	ts := ef.Section(".text")
	if ts == nil {
		return errors.New("no .text section")
	}

	if ts.Addr < minAddr || ts.Addr+uint64(len(text)) > fileEnd {
		return fmt.Errorf(".text %#x+%d outside image [%#x, %#x)",
			ts.Addr, len(text), minAddr, fileEnd)
	}

	copy(image[ts.Addr-minAddr:], text)

	// Size = memory [minAddr, maxEnd), the file only fileEnd-minAddr:
	// the difference (the NOBITS tail) goes into p_memsz
	blob, werr := file.WriteELF(machine, elfFlags(hdr.Machine), minAddr, hdr.Entry, []file.Section{
		*file.NewSection(".image", "", minAddr, 0, maxEnd-minAddr, image),
	})
	if werr != nil {
		return werr
	}

	return os.WriteFile(outPath, blob, 0o755)
}

// allocSection - a section that occupies memory in the image (SHF_ALLOC with
// a non-zero address; metadata tables outside the image have address zero).
func allocSection(s *elf.Section) bool {
	return s.Addr != 0 && s.Flags&elf.SHF_ALLOC != 0
}

// elfMachines - the e_machine values of the originals the writer can rebuild.
var elfMachines = map[elf.Machine]uint16{
	elf.EM_AARCH64:   file.EM_AARCH64,
	elf.EM_RISCV:     file.EM_RISCV,
	elf.EM_LOONGARCH: file.EM_LOONGARCH,
}

// elfFlags — e_flags by machine: the LoongArch kernel requires the ABI
// bits (base + double-float); everything else writes zero.
func elfFlags(machine elf.Machine) uint32 {
	if machine == elf.EM_LOONGARCH {
		return 0x43
	}

	return 0
}
