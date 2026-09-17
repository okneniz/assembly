package arm64

// The Mach-O end of the program chain: assembling with the writer's own
// placement and wrapping as a native image - one call from a built
// program to executable bytes, no address computed anywhere outside the
// file package.

import (
	"fmt"
	"slices"

	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/prog"
)

// MachO - assemble the program with the Mach-O streams policy (the same
// placement the writer uses) and wrap the result as a native arm64
// Mach-O image. entry names the entry symbol. Every label of the program
// becomes a global symbol.
func (b *Binary) MachO(entry string) (*file.MachOImage, error) {
	res := b.AssembleLayout(file.MachoPlaceStreams)
	if len(res.Errs) > 0 {
		return nil, fmt.Errorf("macho: %w", res.Errs[0])
	}

	sections := []file.MachOSection{
		{Segment: "__TEXT", Name: "__text", Data: res.Code, Align: 4},
	}

	if len(res.Data) > 0 {
		sections = append(sections, file.MachOSection{
			Segment: "__DATA", Name: "__data", Data: res.Data, Align: 8,
		})
	}

	if tail := res.DataMem - len(res.Data); tail > 0 {
		sections = append(sections, file.MachOSection{
			Segment: "__DATA", Name: "__bss", Nobits: tail, Align: 8,
		})
	}

	addrs, err := file.MachoPlaceSections(sections)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(res.Syms))
	for name := range res.Syms {
		names = append(names, name)
	}
	slices.Sort(names)

	syms := make([]file.MachOSym, 0, len(names))
	for _, name := range names {
		at, err := machoSectOf(res.Syms[name], sections, addrs)
		if err != nil {
			return nil, err
		}

		syms = append(syms, file.MachOSym{
			Name:    name,
			Section: sections[at].Name,
			Off:     res.Syms[name] - addrs[at],
			Global:  true,
		})
	}

	return file.NewMachOImage(sections, syms, entry)
}

// machoSectOf - the index of the section holding addr.
func machoSectOf(addr uint64, sections []file.MachOSection, addrs []uint64) (int, error) {
	for i := range sections {
		size := uint64(len(sections[i].Data))
		if sections[i].Nobits > 0 {
			size = uint64(sections[i].Nobits)
		}

		if addr >= addrs[i] && addr < addrs[i]+size {
			return i, nil
		}
	}

	return -1, fmt.Errorf("macho: symbol at %#x sits in no section of the program", addr)
}

// compile-time: the streams policy is the shared Place shape.
var _ prog.Place = file.MachoPlaceStreams
