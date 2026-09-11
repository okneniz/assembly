// Package loong64 is the loong64 debug target: the register layout of
// the qemu gdbstub, the executor command line of the virt machine, and
// the listing glue over the arch decoder.
package loong64

import (
	"fmt"

	parsecbytes "github.com/okneniz/parsec/bytes"

	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/debug"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/text"
)

// The gdbstub register numbers of the la64 core layout: $r0-$r31
// (0-31), orig_a0 (32), pc (33); $r3 doubles as sp. flashBase is the
// address the machine resets into - the loader must place the image
// there (as a raw blob: the virt machine's loader device ignores a
// loong ELF).
const (
	pcNum    = 33
	spNum    = 3
	origA0   = 32
	flashMem = 0x1c000000
)

// Target - the loong64 implementation of debug.Target.
type Target struct{}

// NewTarget - the stateless loong64 target.
func NewTarget() Target {
	return Target{}
}

// Arch is "loong64".
func (Target) Arch() string {
	return "loong64"
}

// QemuBinary - the virt machine executor.
func (Target) QemuBinary() string {
	return "qemu-system-loongarch64"
}

// PCNum and SPNum - the gdbstub numbers (pc follows the 32 GPRs).
func (Target) PCNum() int {
	return pcNum
}

func (Target) SPNum() int {
	return spNum
}

// QemuArgs - the virt machine with the raw image placed at the flash
// base the machine resets into (the loader's ELF path is a no-op on
// this machine; the raw blob with an explicit address is the load
// mechanism that works).
func (Target) QemuArgs(imgPath string) []string {
	return []string{
		"-machine", "virt",
		"-device", fmt.Sprintf("loader,file=%s,addr=%#x,cpu-num=0", imgPath, flashMem),
	}
}

// Registers is the ordered core register set: $r0-$r31, orig_a0, pc
// (orig_a0 is part of the stub layout - it keeps the block offsets
// aligned even though the display barely needs it).
func (Target) Registers() []debug.Reg {
	regs := make([]debug.Reg, 0, 34)
	for i := range 32 {
		regs = append(regs, debug.NewReg(fmt.Sprintf("$r%d", i), i, 64))
	}

	return append(
		regs,
		debug.NewReg("orig_a0", origA0, 64),
		debug.NewReg("pc", pcNum, 64),
	)
}

// Disasm - the listing lines of the buffer at addr through the loong64
// decoder (the decode tail falls back to .word lines on its own).
func (Target) Disasm(code []byte, addr uint64) []string {
	instrs, err := arch.MakeDecoder()(parsecbytes.Buffer(code))
	if err != nil {
		return nil
	}

	opts := disasm.NewOptions(text.CodeWord)
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

// InstrLen - a fixed 4 bytes (the LoongArch encoding has no compressed
// form).
func (Target) InstrLen([]byte) int {
	return 4
}
