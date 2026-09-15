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

// mustReg - the pre-minting helper: the names/numbers/widths are
// constants, the constructor error is unreachable.
func mustReg(name string, num, bits int) debug.Reg {
	r, err := debug.NewReg(name, num, bits)
	if err != nil {
		panic(err) // unreachable: the inputs are constants
	}

	return r
}

// Registers is the ordered core register set of the 'g' block prefix:
// x0-x30, sp, pc, cpsr. The vector registers of the full block are not
// described - the display shows the core set.
func (Target) Registers() []debug.Reg {
	regs := make([]debug.Reg, 0, 34)
	for i := range 31 {
		regs = append(regs, mustReg(fmt.Sprintf("x%d", i), i, 64))
	}

	regs = append(
		regs,
		mustReg("sp", spNum, 64),
		mustReg("pc", pcNum, 64),
		mustReg("cpsr", cpsrNum, 32),
	)

	return regs
}
