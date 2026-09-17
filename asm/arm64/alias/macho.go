package alias

// The Mach-O end of the .s path: assembling with the writer's own
// placement (the gas-section policy of the universal Mach-O writer) and
// wrapping as a native image - one call from source text to executable
// bytes, no address computed anywhere outside the file package.

import (
	"fmt"
	"slices"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/file"
)

// MachoLayout is the asm.AssembleLayout policy of the Mach-O writer: the
// gas sections of the source placed exactly where the image puts them.
func MachoLayout(specs []asm.SectionSpec) ([]uint64, error) {
	gas := make([]file.MachoGasSpec, len(specs))
	for i := range specs {
		gas[i] = file.NewMachoGasSpec(specs[i].Name, specs[i].Size, specs[i].Nobits)
	}

	return file.MachoPlaceGas(gas)
}

// MachO assembles the source with MachoLayout and wraps the result as a
// native arm64 MH_EXECUTE (MachOFromResult does the wrapping).
func MachO(src string, entry string) (*file.MachOImage, []asm.AsmError) {
	res, errs := asm.AssembleLayout(src, NewASMBackend(), MachoLayout)
	if len(errs) > 0 {
		return nil, errs
	}

	img, err := MachOFromResult(res, entry)
	if err != nil {
		return nil, []asm.AsmError{asm.NewAsmError(0, 0, err.Error())}
	}

	return img, nil
}

// MachOFromResult wraps an assembled result (alias.AssembleLayout with
// MachoLayout) as a native arm64 MH_EXECUTE: .text into __TEXT,
// .data/.bss (and custom sections) into __DATA, symbols resolved against
// the writer's placement, .global symbols exported. An empty entry picks
// start/_start, and a source with neither enters at the first
// instruction of .text (as _main).
func MachOFromResult(res *asm.Result, entry string) (*file.MachOImage, error) {
	sections := make([]file.MachOSection, len(res.Sections))
	for i, s := range res.Sections {
		reserve := 0
		if s.Nobits {
			reserve = s.Size
		}

		sections[i] = file.MachoGasSection(s.Name, s.Data, reserve)
	}

	globals := make(map[string]bool, len(res.Globals))
	for _, name := range res.Globals {
		globals[name] = true
	}

	names := make([]string, 0, len(res.Symbols))
	for name := range res.Symbols {
		names = append(names, name)
	}
	slices.Sort(names)

	syms := make([]file.MachOSym, 0, len(names))
	for _, name := range names {
		addr := res.Symbols[name]

		at, err := machoSectAt(addr, res.Sections)
		if err != nil {
			return nil, err
		}

		syms = append(syms, file.MachOSym{
			Name:    name,
			Section: sections[at].Name,
			Off:     addr - res.Sections[at].Addr,
			Global:  globals[name],
		})
	}

	if entry == "" {
		for _, name := range []string{"start", "_start"} {
			if _, ok := res.Symbols[name]; ok {
				entry = name
				break
			}
		}
	}

	if entry == "" {
		entry = "_main"
		syms = append(syms, file.MachOSym{
			Name:    "_main",
			Section: sections[0].Name,
			Global:  true,
		})
	}

	return file.NewMachOImage(sections, syms, entry)
}

// machoSectAt - the index of the result section holding addr.
func machoSectAt(addr uint64, secs []asm.Section) (int, error) {
	for i, s := range secs {
		if addr >= s.Addr && addr < s.Addr+uint64(s.Size) {
			return i, nil
		}
	}

	return -1, fmt.Errorf(
		"macho: a symbol at %#x sits in no section of the source", addr,
	)
}
