package file

// Adapter of the file/macho subpackage to the neutral FileFormat contract.

import (
	"fmt"

	"github.com/okneniz/assembly/file/macho"
)

// machoFormat wraps *macho.File.
type machoFormat struct {
	f *macho.File
}

func (x machoFormat) Name() string {
	return "Mach-O"
}

func (x machoFormat) ArchKind() ArchKind {
	if x.f.Header().CpuType == macho.CPU_TYPE_ARM64 {
		return ArchARM64
	}

	return ArchUnknown
}

func (x machoFormat) Sections() ([]Section, error) {
	src := x.f.Sections()
	out := make([]Section, 0, len(src))
	for _, s := range src {
		out = append(out, *NewSection(s.SectName, s.SegName, s.Addr, uint64(s.Offset), s.Size, nil))
	}

	return out, nil
}

func (x machoFormat) Section(name string) (*Section, error) {
	s := x.f.Section(name)
	if s == nil {
		return nil, fmt.Errorf("assembly/macho: section %q not found", name)
	}

	data, err := s.Data()
	if err != nil {
		return nil, err
	}

	return NewSection(s.SectName, s.SegName, s.Addr, uint64(s.Offset), s.Size, data), nil
}

// CodeSection returns the code section (__text of segment __TEXT).
func (x machoFormat) CodeSection() (*Section, error) {
	return x.Section("__text")
}
