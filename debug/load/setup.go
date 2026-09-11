package load

import (
	"fmt"
	"strings"

	"github.com/okneniz/assembly/asm"
	"github.com/okneniz/assembly/asm/arm64/alias"
	lpseudo "github.com/okneniz/assembly/asm/loong64/pseudo"
	rpseudo "github.com/okneniz/assembly/asm/riscv/pseudo"
	"github.com/okneniz/assembly/debug"
	darm64 "github.com/okneniz/assembly/debug/arm64"
	dloong "github.com/okneniz/assembly/debug/loong64"
	driscv "github.com/okneniz/assembly/debug/riscv"
	"github.com/okneniz/assembly/file"
)

// Setup is the per-arch configuration of the loading path: the debug
// target, the text assembler, the ELF machine/flags, the default
// section base of the arch's qemu machine convention, and the
// raw-image flag (the loong machine loads a raw blob at its flash
// base, not an ELF).
type Setup struct {
	Tgt      debug.Target
	Assemble func(string, uint64) (*asm.Result, []asm.AsmError)
	Machine  uint16
	Flags    uint32
	Base     uint64
	Raw      bool
}

// SetupFor resolves the arch name (arm64|riscv64|loong64) into its
// configuration.
func SetupFor(name string) (Setup, error) {
	switch strings.ToLower(name) {
	case "arm64", "aarch64":
		return Setup{
			Tgt:      darm64.NewTarget(),
			Assemble: alias.Assemble,
			Machine:  file.EM_AARCH64,
			Flags:    0,
			Base:     0x40100000,
		}, nil
	case "riscv64", "riscv":
		return Setup{
			Tgt:      driscv.NewTarget(),
			Assemble: rpseudo.Assemble,
			Machine:  file.EM_RISCV,
			Flags:    0,
			Base:     0x80000000,
		}, nil
	case "loong64", "loongarch64":
		return Setup{
			Tgt:      dloong.NewTarget(),
			Assemble: lpseudo.Assemble,
			Machine:  file.EM_LOONGARCH,
			Flags:    0x43,
			Base:     0x1c000000,
			Raw:      true,
		}, nil
	}

	return Setup{}, fmt.Errorf(
		"assembly/load: unknown arch %q (want arm64, riscv64 or loong64)",
		name,
	)
}
