// Package arm64 - the arm64 end of the in-memory linker: the arch
// assembler injected into the arch-neutral driver, the resolved output
// wrapped by the platform writers - a native Mach-O image for the host,
// an executable ELF for Linux (qemu, the docker gates).
package arm64

import (
	"fmt"

	"github.com/okneniz/assembly/asm/arm64/alias"
	"github.com/okneniz/assembly/file"
	"github.com/okneniz/assembly/link"
	unitarm64 "github.com/okneniz/assembly/unit/arm64"
)

// Macho links the sources into a native arm64 Mach-O image: the
// writer's own placement resolves the streams, every label of the
// program becomes a global symbol, entry names the entry symbol (""
// picks start/_start).
func Macho(sources []link.Source, entry string) (*file.MachOImage, error) {
	f := link.Link(deps(), sources, entry, file.MachoPlaceStreams)
	if len(f.Errs) > 0 {
		return nil, fmt.Errorf("link: %w", f.Errs[0])
	}

	name, ok := link.EntryOf(f, entry)
	if !ok {
		return nil, fmt.Errorf("link: no entry symbol (want %q, start or _start)", entry)
	}

	return unitarm64.MachoOf(f, name)
}

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

	return file.WriteELF(file.EM_AARCH64, 0, base, f.Syms[name], sections)
}

// deps is the arm64 injection of the link driver.
func deps() link.Deps {
	return link.Deps{Parse: alias.ParseSourceUnit}
}
