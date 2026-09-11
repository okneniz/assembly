// Package riscv is the riscv64 debug target: the register layout of
// the qemu gdbstub, the executor command line of the virt machine, and
// the listing glue over the arch decoder (including the compressed
// instruction lengths of RVC).
package riscv

import (
	"fmt"

	parsecbytes "github.com/okneniz/parsec/bytes"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/debug"
	"github.com/okneniz/assembly/disasm"
	"github.com/okneniz/assembly/text"
)

// The gdbstub register numbers of the rv64 core layout: x0-x31 (0-31),
// pc (32); x2 doubles as sp.
const (
	pcNum = 32
	spNum = 2
)

// Target - the riscv64 implementation of debug.Target.
type Target struct{}

// NewTarget - the stateless riscv64 target.
func NewTarget() Target {
	return Target{}
}

// Arch is "riscv64".
func (Target) Arch() string {
	return "riscv64"
}

// QemuBinary - the virt machine executor.
func (Target) QemuBinary() string {
	return "qemu-system-riscv64"
}

// PCNum and SPNum - the gdbstub numbers (pc follows the 32 GPRs).
func (Target) PCNum() int {
	return pcNum
}

func (Target) SPNum() int {
	return spNum
}

// QemuArgs - the virt machine with the ELF as the kernel: qemu maps it
// at its p_vaddr and starts the hart at the entry (no firmware).
func (Target) QemuArgs(imgPath string) []string {
	return []string{
		"-machine", "virt",
		"-bios", "none",
		"-kernel", imgPath,
	}
}

// Registers is the ordered core register set: x0-x31, pc (the plain
// stub names - the ABI aliases live in the listing, not the dump).
func (Target) Registers() []debug.Reg {
	regs := make([]debug.Reg, 0, 33)
	for i := range 32 {
		regs = append(regs, debug.NewReg(fmt.Sprintf("x%d", i), i, 64))
	}

	return append(regs, debug.NewReg("pc", pcNum, 64))
}

// Disasm - the listing lines of the buffer at addr through the riscv
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

// InstrLen - 2 for a compressed (RVC) instruction, 4 otherwise: the
// low two bits 00-10 mark the compressed space, 11 the 32-bit one. A
// nil or short buffer yields the arch constant 4 (the breakpoint
// kind).
func (Target) InstrLen(code []byte) int {
	if len(code) < 2 {
		return 4
	}

	if code[0]&0x3 != 0x3 {
		return 2
	}

	return 4
}
