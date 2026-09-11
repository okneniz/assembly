package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/debug"
)

func TestRegisters(t *testing.T) {
	regs := NewTarget().Registers()
	require.Len(t, regs, 34)

	// the 'g' block prefix: x0-x30, sp, pc, cpsr - the numbers are the
	// RSP register numbers, the order is the block order
	require.Equal(t, debug.NewReg("x0", 0, 64), regs[0])
	require.Equal(t, debug.NewReg("x30", 30, 64), regs[30])
	require.Equal(t, debug.NewReg("sp", 31, 64), regs[31])
	require.Equal(t, debug.NewReg("pc", 32, 64), regs[32])
	require.Equal(t, debug.NewReg("cpsr", 33, 32), regs[33])
}

func TestNums(t *testing.T) {
	tgt := NewTarget()
	require.Equal(t, 32, tgt.PCNum())
	require.Equal(t, 31, tgt.SPNum())
}

func TestQemuArgs(t *testing.T) {
	args := NewTarget().QemuArgs("/tmp/img.elf")
	require.Equal(t, []string{
		"-machine", "virt",
		"-cpu", "cortex-a53",
		"-device", "loader,file=/tmp/img.elf,cpu-num=0",
	}, args)
}

func TestInstrLen(t *testing.T) {
	require.Equal(t, 4, NewTarget().InstrLen(nil))
	require.Equal(t, 4, NewTarget().InstrLen([]byte{1, 2, 3}))
}

func TestDisasm(t *testing.T) {
	// movz x0, #1 (0xd2800020) at 0x401000: one listing line with the
	// address, the code, and the decoded text
	lines := NewTarget().Disasm([]byte{0x20, 0x00, 0x80, 0xd2}, 0x401000)
	require.Len(t, lines, 1)
	require.Contains(t, lines[0], "401000:")
	require.Contains(t, lines[0], "20 00 80 d2")
	require.Contains(t, lines[0], "mov x0, #0x1")
}
