package loong64

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/debug"
)

func TestRegisters(t *testing.T) {
	regs := NewTarget().Registers()
	require.Len(t, regs, 34)
	require.Equal(t, debug.NewReg("$r0", 0, 64), regs[0])
	require.Equal(t, debug.NewReg("$r31", 31, 64), regs[31])
	require.Equal(t, debug.NewReg("orig_a0", 32, 64), regs[32])
	require.Equal(t, debug.NewReg("pc", 33, 64), regs[33])
}

func TestNums(t *testing.T) {
	tgt := NewTarget()
	require.Equal(t, 33, tgt.PCNum())
	require.Equal(t, 3, tgt.SPNum())
}

func TestQemuArgs(t *testing.T) {
	require.Equal(t, []string{
		"-machine", "virt",
		"-device", "loader,file=/tmp/img.bin,addr=0x1c000000,cpu-num=0",
	}, NewTarget().QemuArgs("/tmp/img.bin"))
}

func TestInstrLen(t *testing.T) {
	require.Equal(t, 4, NewTarget().InstrLen(nil))
	require.Equal(t, 4, NewTarget().InstrLen([]byte{1, 2, 3}))
}

func TestDisasm(t *testing.T) {
	// ori $t0, $zero, 1 at 0x1c000000 (word-style code column), the
	// instruction built by the arch builder - no hand-encoded bytes
	b := arch.Builder{}
	rd, err := arch.R(12)
	require.NoError(t, err)
	rj, err := arch.R(0)
	require.NoError(t, err)

	imm, err := b.UImm12(1)
	require.NoError(t, err)

	var code bytes.Buffer
	_, err = b.Ori(rd, rj, imm).Encode(&code)
	require.NoError(t, err)

	lines := NewTarget().Disasm(code.Bytes(), 0x1c000000)
	require.Len(t, lines, 1)
	require.Contains(t, lines[0], "1c000000:")
	require.Contains(t, lines[0], "ori")
}
