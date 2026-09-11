package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/debug"
)

func TestRegisters(t *testing.T) {
	regs := NewTarget().Registers()
	require.Len(t, regs, 33)
	require.Equal(t, debug.NewReg("x0", 0, 64), regs[0])
	require.Equal(t, debug.NewReg("x31", 31, 64), regs[31])
	require.Equal(t, debug.NewReg("pc", 32, 64), regs[32])
}

func TestNums(t *testing.T) {
	tgt := NewTarget()
	require.Equal(t, 32, tgt.PCNum())
	require.Equal(t, 2, tgt.SPNum())
}

func TestQemuArgs(t *testing.T) {
	require.Equal(t, []string{
		"-machine", "virt",
		"-bios", "none",
		"-kernel", "/tmp/img.elf",
	}, NewTarget().QemuArgs("/tmp/img.elf"))
}

func TestInstrLen(t *testing.T) {
	cases := []struct {
		name string
		code []byte
		want int
	}{
		{"nil buffer", nil, 4},
		{"one byte", []byte{0x01}, 4},
		{"32-bit (low bits 11)", []byte{0x93, 0x02, 0x10, 0x00}, 4},
		{"compressed c.addi (low bits 00)", []byte{0x81, 0x02}, 2},
		{"compressed (low bits 01)", []byte{0x01, 0x45}, 2},
		{"compressed (low bits 10)", []byte{0xfe, 0x4f}, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.want, NewTarget().InstrLen(c.code))
		})
	}
}

func TestDisasm(t *testing.T) {
	// addi x5, x0, 1 at 0x80000000 (word-style code column; the
	// decoder prints the li alias of the addi form)
	lines := NewTarget().Disasm([]byte{0x93, 0x02, 0x10, 0x00}, 0x80000000)
	require.Len(t, lines, 1)
	require.Contains(t, lines[0], "80000000:")
	require.Contains(t, lines[0], "00100293")
	require.Contains(t, lines[0], "li t0, 0x1")
}
