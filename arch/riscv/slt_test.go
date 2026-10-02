package riscv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSltCtor(t *testing.T) {
	for _, c := range []struct {
		name  string
		instr Instr
		word  uint32
	}{
		{"slt a0, a1, a2", New().Slt(xreg(t, 10), xreg(t, 11), xreg(t, 12)), 0x00c5a533},
		{"slt zero, t0, t1", New().Slt(xreg(t, 0), xreg(t, 5), xreg(t, 6)), 0x0062a033},
		{"slt t6, ra, sp", New().Slt(xreg(t, 31), xreg(t, 1), xreg(t, 2)), 0x0020afb3},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.word, ctorWord(t, c.instr))
		})
	}
}

// TestSltPseudoCanon - the llvm-mc print canon of the zero-operand
// shapes: sltz (rs2 = x0) and sgtz (rs1 = x0).
func TestSltPseudoCanon(t *testing.T) {
	for _, c := range []struct {
		name string
		in   Instr
		text string
	}{
		{"sltz a0, a1", New().Slt(xreg(t, 10), xreg(t, 11), xreg(t, 0)), "sltz a0, a1"},
		{"sgtz a0, a1", New().Slt(xreg(t, 10), xreg(t, 0), xreg(t, 11)), "sgtz a0, a1"},
		{"sltz zero, zero", New().Slt(xreg(t, 0), xreg(t, 0), xreg(t, 0)), "sltz zero, zero"},
	} {
		t.Run(c.name, func(t *testing.T) {
			require.Equal(t, c.text, c.in.ObjDump(nil))
		})
	}
}
