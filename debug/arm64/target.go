// Package arm64 is the arm64 debug target: the register layout of the
// qemu gdbstub, the executor command line of the virt machine, and the
// listing glue over the arch decoder. A value type - the stateless
// data half of the debugger.
package arm64

import (
	"fmt"

	parsecbytes "github.com/okneniz/parsec/bytes"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/text"
)

// Target - the arm64 implementation of debug.Target.
type Target struct{}

// NewTarget - the stateless arm64 target.
func NewTarget() Target {
	return Target{}
}

// Arch is "arm64".
func (Target) Arch() string {
	return "arm64"
}

// QemuBinary - the virt machine executor.
func (Target) QemuBinary() string {
	return "qemu-system-aarch64"
}

// PCNum and SPNum - the gdbstub numbers: pc follows the 31 GPRs and sp.
func (Target) PCNum() int {
	return pcNum
}

func (Target) SPNum() int {
	return spNum
}

// QemuArgs - the virt machine with the image loaded by the loader
// device: it maps the ELF segments and starts cpu 0 at the entry.
func (Target) QemuArgs(imgPath string) []string {
	return []string{
		"-machine", "virt",
		"-cpu", "cortex-a53",
		"-device", fmt.Sprintf("loader,file=%s,cpu-num=0", imgPath),
	}
}

// Disasm - the listing lines of the buffer at addr through the arm64
// decoder (the decode tail falls back to .word lines on its own).
func (Target) Disasm(code []byte, addr uint64) []string {
	instrs, err := arch.MakeDecoder()(parsecbytes.Buffer(code))
	if err != nil {
		return nil
	}

	opts := disasm.NewOptions(text.CodeBytes)
	out := make([]string, 0, len(instrs))
	off := 0
	for _, in := range instrs {
		if off < len(code) {
			out = append(out, disasm.Line(addr+uint64(off), code[off:], in, opts))
		}

		off += in.Len()
	}

	return out
}

// InstrLen - a fixed 4 bytes (no arm64 compressed form).
func (Target) InstrLen([]byte) int {
	return 4
}
