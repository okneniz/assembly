package file

// Adapter of the file/elf subpackage to the neutral FileFormat contract.

import (
	"fmt"

	"github.com/okneniz/assembly/file/elf"
)

// elfFormat wraps *elf.File.
type elfFormat struct {
	f *elf.File
}

func (x elfFormat) Name() string {
	return "ELF"
}

func (x elfFormat) ArchKind() ArchKind {
	switch x.f.Header().Machine {
	case elf.EM_AARCH64:
		return ArchARM64
	case elf.EM_RISCV:
		return ArchRISCV64
	case elf.EM_LOONGARCH:
		return ArchLOONGARCH64
	default:
		return ArchUnknown
	}
}

func (x elfFormat) Sections() ([]Section, error) {
	src := x.f.Sections()
	out := make([]Section, 0, len(src))
	for _, s := range src {
		out = append(out, *NewSection(s.Name, "", s.Addr, s.Off, s.Size, nil))
	}

	return out, nil
}

func (x elfFormat) Section(name string) (*Section, error) {
	s := x.f.Section(name)
	if s == nil {
		return nil, fmt.Errorf("assembly/elf: section %q not found", name)
	}

	data, err := s.Data()
	if err != nil {
		return nil, err
	}

	return NewSection(s.Name, "", s.Addr, s.Off, s.Size, data), nil
}

// CodeSection returns the code section (.text).
func (x elfFormat) CodeSection() (*Section, error) {
	return x.Section(".text")
}
