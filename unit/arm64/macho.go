package arm64

// The Mach-O end of the unit output: resolving with the writer's own
// placement and wrapping as a native image - one call from a Unit to
// executable bytes, no address computed anywhere outside the file
// package.

import (
	"fmt"
	"slices"

	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/unit"
)

// Macho - resolve the output with the Mach-O streams policy (the same
// placement the writer uses) and wrap it as a native arm64 Mach-O image.
// entry names the entry symbol; every label of the program becomes a
// global symbol.
func Macho(u *unit.Unit, entry string) (*file.MachOImage, error) {
	f := u.Resolve(file.MachoPlaceStreams)
	if len(f.Errs) > 0 {
		return nil, fmt.Errorf("macho: %w", f.Errs[0])
	}

	code, codeErr := f.EncodeText()
	if codeErr != nil {
		return nil, fmt.Errorf("macho: %w", codeErr)
	}

	data, dataErr := f.EncodeData()
	if dataErr != nil {
		return nil, fmt.Errorf("macho: %w", dataErr)
	}

	sections := []file.MachOSection{
		file.NewMachOSection("__TEXT", "__text", code, 0, 4),
	}

	if len(data) > 0 {
		sections = append(sections, file.NewMachOSection(
			"__DATA", "__data", data, 0, 8,
		))
	}

	if tail := f.DataMem - len(data); tail > 0 {
		sections = append(sections, file.NewMachOSection(
			"__DATA", "__bss", nil, tail, 8,
		))
	}

	addrs, err := file.MachoPlaceSections(sections)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(f.Syms))
	for name := range f.Syms {
		names = append(names, name)
	}
	slices.Sort(names)

	syms := make([]file.MachOSym, 0, len(names))
	for _, name := range names {
		at, err := machoSectOf(f.Syms[name], sections, addrs)
		if err != nil {
			return nil, err
		}

		syms = append(syms, file.NewMachOSym(
			name, sections[at].Name, f.Syms[name]-addrs[at], true,
		))
	}

	return file.NewMachOImage(sections, syms, entry)
}

// machoSectOf is the index of the section holding addr. A label may also
// sit exactly at the end of its stream (the fixup tables of a compiler
// close on one, the linker's _end is the same shape): containment wins
// first, a boundary address falls to the section it ends - at offset
// == size.
func machoSectOf(addr uint64, sections []file.MachOSection, addrs []uint64) (int, error) {
	end := -1
	for i := range sections {
		size := uint64(len(sections[i].Data))
		if sections[i].Nobits > 0 {
			size = uint64(sections[i].Nobits)
		}

		if addr >= addrs[i] && addr < addrs[i]+size {
			return i, nil
		}

		if addr == addrs[i]+size {
			end = i
		}
	}

	if end >= 0 {
		return end, nil
	}

	return -1, fmt.Errorf("macho: symbol at %#x sits in no section of the program", addr)
}

// compile-time: the streams policy is the shared Place shape.
var _ unit.Place = file.MachoPlaceStreams
