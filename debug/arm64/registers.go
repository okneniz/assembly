package arm64

import (
	"fmt"

	"github.com/okneniz/assembly/debug"
)

// The gdbstub register numbers of the aarch64 core layout: x0-x30
// (0-30), sp (31), pc (32), cpsr/pstate (33, 32-bit).
const (
	pcNum   = 32
	spNum   = 31
	cpsrNum = 33
)

// Registers is the ordered core register set of the 'g' block prefix:
// x0-x30, sp, pc, cpsr. The vector registers of the full block are not
// described - the display shows the core set.
func (Target) Registers() []debug.Reg {
	regs := make([]debug.Reg, 0, 34)
	for i := range 31 {
		regs = append(regs, debug.NewReg(fmt.Sprintf("x%d", i), i, 64))
	}

	regs = append(
		regs,
		debug.NewReg("sp", spNum, 64),
		debug.NewReg("pc", pcNum, 64),
		debug.NewReg("cpsr", cpsrNum, 32),
	)

	return regs
}
