// Package loong64 - the LoongArch64 end of the in-memory linker: the
// arch assembler injected into the arch-neutral driver, the resolved
// output wrapped by the platform writer. The target is Linux: an
// executable ELF (qemu, the future loong container) - there is no
// Mach-O on this architecture. The ABI flag is the double-float base
// (0x43), the one the emitter of the CLI carries.
package loong64

import (
	"fmt"

	"github.com/okneniz/assembly/asm/loong64/pseudo"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/link"
)

// ELF links the sources into an executable ELF64 at base: a flat image
// - the text at base, the data past it (the emitter lays the sections
// back to back, so the data base is base+len(text), no gap), the bss as
// the data section's NOBITS tail. entry names the entry symbol (""
// picks start/_start).
func ELF(sources []link.Source, entry string, base uint64) ([]byte, error) {
	place := func(text, data, dataMem int) (uint64, uint64) {
		return base, base + uint64(text)
	}

	f := link.Link(deps(), sources, entry, place)
	if len(f.Errs) > 0 {
		return nil, fmt.Errorf("link: %w", f.Errs[0])
	}

	name, ok := link.EntryOf(f, entry)
	if !ok {
		return nil, fmt.Errorf("link: no entry symbol (want %q, start or _start)", entry)
	}

	code, err := f.EncodeText()
	if err != nil {
		return nil, fmt.Errorf("elf: %w", err)
	}

	data, err := f.EncodeData()
	if err != nil {
		return nil, fmt.Errorf("elf: %w", err)
	}

	sections := []file.Section{
		*file.NewSection("text", "", base, 0, uint64(len(code)), code),
		*file.NewSection("data", "", base+uint64(len(code)), 0, uint64(f.DataMem), data),
	}

	return file.WriteELF(file.EM_LOONGARCH, 0x43, base, f.Syms[name], sections)
}

// deps is the LoongArch64 injection of the link driver.
func deps() link.Deps {
	return link.Deps{Parse: pseudo.ParseSourceUnit}
}
