package riscv

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// --- operand generators --------------------------------------------------

// xreg — a register by number for tests: an X(n) error fails the test.
func xreg(t *testing.T, n int) riscv.Reg {
	t.Helper()
	r, err := riscv.X(n)
	require.NoError(t, err)
	return r
}

// TestRegGenProperty — samples in x0..x31, shrink to zero.
func TestRegGenProperty(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 500, Reg(rnd), func(r riscv.Reg) bool {
		return r.Num() <= 31
	})
}

func TestRegGenDistribution(t *testing.T) {
	rnd := arb.Rnd(42)
	seen := map[string]bool{}
	for range 1000 {
		seen[ohsnap.First(Reg(rnd).Generate()).String()] = true
	}

	for _, want := range []string{"zero", "ra", "sp", "t0", "a0", "t6"} {
		require.True(t, seen[want], "%s was not seen in 1000 samples", want)
	}
}

func TestRegGenShrink(t *testing.T) {
	got := regShrunk(xreg(t, 9))
	require.Len(t, got, 2, "regShrunk(x9) = %v", got)
	require.Equal(t, riscv.Zero, got[0])
	require.Equal(t, xreg(t, 4), got[1])
	got = regShrunk(riscv.Zero)
	require.Empty(t, got, "regShrunk(zero) — shrinker self-loop")
	got = regShrunk(xreg(t, 1))
	require.Len(t, got, 1, "regShrunk(x1) = %v", got)
	require.Equal(t, riscv.Zero, got[0])
}

// TestImmGenProperty — signed immediates in range, shrink in range.
func TestImmGenProperty(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 300, Imm12(rnd), func(v riscv.Imm12) bool {
		n, err := immValue(v)
		if err != nil {
			t.Errorf("immValue(%v): %v", v, err)
			return false
		}

		return n >= -2048 && n <= 2047
	})
	ohsnap.Check(t, 300, Imm20(rnd), func(v riscv.Imm20) bool {
		n, err := immValue(v)
		if err != nil {
			t.Errorf("immValue(%v): %v", v, err)
			return false
		}

		return n >= 0 && n <= 0xfffff
	})
	ohsnap.Check(t, 300, Off(rnd), func(v riscv.Off) bool {
		n, err := immValue(v)
		if err != nil {
			t.Errorf("immValue(%v): %v", v, err)
			return false
		}

		return n >= -2048 && n <= 2047
	})
	lo, err := riscv.NewImm12(-2048)
	require.NoError(t, err)
	for s := range Imm12(rnd).Shrink(lo) {
		n, err := immValue(s)
		require.NoError(t, err)
		require.GreaterOrEqual(t, n, int64(-2048), "Imm12 shrink out of range: %v", s)
		require.LessOrEqual(t, n, int64(2047), "Imm12 shrink out of range: %v", s)
	}

	// The range boundaries come first (shrink.Boundaries), halving toward
	// zero follows and ends with zero.
	mid, err := riscv.NewImm12(100)
	require.NoError(t, err)
	cs := slices.Collect(Imm12(rnd).Shrink(mid))
	require.NotEmpty(t, cs)
	first, err := immValue(cs[0])
	require.NoError(t, err)
	second, err := immValue(cs[1])
	require.NoError(t, err)
	require.Equal(t, int64(-2048), first, "Imm12 shrink candidates: %v", cs)
	require.Equal(t, int64(2047), second, "Imm12 shrink candidates: %v", cs)
	last, err := immValue(cs[len(cs)-1])
	require.NoError(t, err)
	require.Equal(t, int64(0), last, "Imm12 shrink candidates: %v", cs)
}

// TestImmValue — parsing String() (hex without #, negatives with a minus).
func TestImmValue(t *testing.T) {
	neg, err := riscv.NewImm12(-4)
	require.NoError(t, err)
	n, err := immValue(neg)
	require.NoError(t, err)
	require.Equal(t, int64(-4), n)
	big, err := riscv.NewImm20(0x12345)
	require.NoError(t, err)
	n, err = immValue(big)
	require.NoError(t, err)
	require.Equal(t, int64(0x12345), n)
}

// --- instruction generators ------------------------------------------------

// instrCase — one instruction: name + generator.
type instrCase struct {
	name string
	gen  func() riscv.Instr
}

func newInstrCase(name string, gen func() riscv.Instr) instrCase {
	return instrCase{
		name: name,
		gen:  gen,
	}
}

func instrCases(rnd *rand.Rand) []instrCase {
	return []instrCase{
		newInstrCase("Add", func() riscv.Instr {
			return ohsnap.First(Add(rnd).Generate()).Instr()
		}),
		newInstrCase("Sub", func() riscv.Instr {
			return ohsnap.First(Sub(rnd).Generate()).Instr()
		}),
		newInstrCase("Addi", func() riscv.Instr {
			return ohsnap.First(Addi(rnd).Generate()).Instr()
		}),
		newInstrCase("Lui", func() riscv.Instr {
			return ohsnap.First(Lui(rnd).Generate()).Instr()
		}),
		newInstrCase("Lw", func() riscv.Instr {
			return ohsnap.First(Lw(rnd).Generate()).Instr()
		}),
		newInstrCase("Ld", func() riscv.Instr {
			return ohsnap.First(Ld(rnd).Generate()).Instr()
		}),
		newInstrCase("Sw", func() riscv.Instr {
			return ohsnap.First(Sw(rnd).Generate()).Instr()
		}),
		newInstrCase("Sd", func() riscv.Instr {
			return ohsnap.First(Sd(rnd).Generate()).Instr()
		}),
	}
}

// TestInstrGenValid — a generated instruction always encodes
// (2 or 4 bytes — RVC compression is legal and consistent with the decoder).
func TestInstrGenValid(t *testing.T) {
	rnd := arb.Rnd(42)
	for _, c := range instrCases(rnd) {
		var buf bytes.Buffer
		for range 300 {
			buf.Reset()
			in := c.gen()
			_, err := in.Encode(&buf, 0x1000, riscv.EncOpts{})
			require.NoError(t, err, "%s (%s)", c.name, in.ObjDump(disasm.DefaultViewCtx()))
			require.Contains(
				t,
				[]int{2, 4},
				buf.Len(),
				"%s: %d bytes",
				c.name,
				buf.Len(),
			)
		}
	}
}

// TestInstrGenText — String() of the parameters = ObjDump of the instruction.
func TestInstrGenText(t *testing.T) {
	rnd := arb.Rnd(42)
	got := ohsnap.First(Add(rnd).Generate()).String()
	require.Equal(t, "add ", got[:4], "AddParams.String() = %q", got)
	got = ohsnap.First(Lui(rnd).Generate()).String()
	require.Equal(t, "lui ", got[:4], "LuiParams.String() = %q", got)
}

// TestInstrGenShrinkValid — shrink candidates are valid instructions.
func TestInstrGenShrinkValid(t *testing.T) {
	rnd := arb.Rnd(9)
	addi := Addi(rnd)
	for range 200 {
		p := ohsnap.First(addi.Generate())
		for s := range addi.Shrink(p) {
			s.Instr() // must not panic
		}
	}
}
