package arm64

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// --- operand generators ----------------------------------------------

// mustX/mustW — registers by number in tests: a constructor error here
// means a broken test (the numbers are fixed and known to be in 0..30).
func mustX(t *testing.T, n int) arm64.Reg {
	t.Helper()
	r, err := arm64.X(n)
	require.NoError(t, err, "arm64.X(%d)", n)
	return r
}

func mustW(t *testing.T, n int) arm64.Reg {
	t.Helper()
	r, err := arm64.W(n)
	require.NoError(t, err, "arm64.W(%d)", n)
	return r
}

// mustImmValue — the numeric value of an immediate: a String() parse error
// of our own type is a bug of the type, not of the test.
func mustImmValue(t *testing.T, v any) int64 {
	t.Helper()
	n, err := immValue(v)
	require.NoError(t, err, "immValue(%v)", v)
	return n
}

// TestRegGenProperty — every sample is valid, the name is parsed by arm,
// shrinking preserves the width.
func TestRegGenProperty(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 500, Reg(rnd), func(r arm64.Reg) bool {
		return r.Is64() == strings.HasPrefix(r.String(), "x") ||
			r.String() == "sp" || r.String() == "wsp"
	})
}

func TestRegGenDistribution(t *testing.T) {
	rnd := arb.Rnd(42)
	seen := map[string]bool{}
	for range 1000 {
		seen[ohsnap.First(Reg(rnd).Generate()).String()] = true
	}

	for _, want := range []string{"x0", "w0", "xzr", "wzr", "sp", "wsp", "x30"} {
		require.True(t, seen[want], "%s was not seen in 1000 samples", want)
	}
}

func TestRegGenShrink(t *testing.T) {
	for _, r := range []arm64.Reg{mustX(t, 9), mustW(t, 7), arm64.SP, arm64.XZR} {
		for _, s := range regShrunk(r) {
			require.Equal(t, r.Is64(), s.Is64(), "shrink(%s) = %s changed the width", r, s)
		}
	}

	got := regShrunk(mustX(t, 9))
	require.Len(t, got, 2, "regShrunk(x9) = %v", got)
	require.Equal(t, mustX(t, 0), got[0], "regShrunk(x9) = %v", got)
	require.Equal(t, mustX(t, 4), got[1], "regShrunk(x9) = %v", got)
	require.Empty(t, regShrunk(mustX(t, 0)), "regShrunk(x0) — shrinker self-loop")
}

// TestImmGenProperty — values and all shrink candidates are in range.
func TestImmGenProperty(t *testing.T) {
	checks := []struct {
		name string
		gen  func(rnd *rand.Rand) ohsnap.Arbitrary[int64]
		max  int64
	}{
		{
			"Imm12",
			func(rnd *rand.Rand) ohsnap.Arbitrary[int64] {
				return immArb[int64]{
					rnd: rnd,
					max: 0xfff,
					make: func(v int64) (int64, error) {
						return v, nil
					},
				}
			},
			0xfff,
		},
		{
			"Imm16",
			func(rnd *rand.Rand) ohsnap.Arbitrary[int64] {
				return immArb[int64]{
					rnd: rnd,
					max: 0xffff,
					make: func(v int64) (int64, error) {
						return v, nil
					},
				}
			},
			0xffff,
		},
	}
	for _, c := range checks {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, c.gen(rnd), func(v int64) bool {
			return v >= 0 && v <= c.max
		})
	}

	// Shrink of typed immediates by halving toward zero.
	rnd := arb.Rnd(1)
	i12 := Imm12(rnd)
	i16 := Imm16(rnd)
	for s := range i12.Shrink(ohsnap.First(i12.Generate())) {
		n := mustImmValue(t, s)
		require.True(t, n >= 0 && n <= 0xfff, "Imm12 shrink out of range: %v", s)
	}

	for s := range i16.Shrink(ohsnap.First(i16.Generate())) {
		n := mustImmValue(t, s)
		require.True(t, n >= 0 && n <= 0xffff, "Imm16 shrink out of range: %v", s)
	}

	n := mustImmValue(t, ohsnap.First(Imm6(rnd).Generate()))
	require.GreaterOrEqual(t, n, int64(0), "Imm6: immValue failed to parse String()")
}

func TestImmGenBounds(t *testing.T) {
	// With a fixed seed the bounds are reachable.
	rnd := arb.Rnd(7)
	min12, max12 := int64(1<<62), int64(-1)
	for range 2000 {
		n := mustImmValue(t, ohsnap.First(Imm12(rnd).Generate()))
		min12 = min(min12, n)
		max12 = max(max12, n)
	}

	require.True(
		t,
		min12 <= 0x10 && max12 >= 0xff0,
		"Imm12: min=%#x max=%#x — bounds not covered",
		min12,
		max12,
	)
}

// TestEnumGen — enums from the set, shrink to the first value.
func TestEnumGen(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 200, Shift(rnd), func(s arm64.Shift) bool {
		return s >= arm64.LSL && s <= arm64.ROR
	})
	ohsnap.Check(t, 200, Hw(rnd), func(h arm64.Hw) bool {
		return h >= arm64.Hw0 && h <= arm64.Hw3
	})
	ohsnap.Check(t, 200, Sh12(rnd), func(s arm64.Sh12) bool {
		return s == arm64.NoSh12 || s == arm64.LSL12
	})
	shifts := slices.Collect(Shift(rnd).Shrink(arm64.ROR))
	require.Len(t, shifts, 1, "Shift.Shrink(ROR) = %v", shifts)
	require.Equal(t, arm64.LSL, shifts[0], "Shift.Shrink(ROR) = %v", shifts)
	sh12s := slices.Collect(Sh12(rnd).Shrink(arm64.LSL12))
	require.Len(t, sh12s, 1, "Sh12.Shrink(LSL12) = %v", sh12s)
	require.Equal(t, arm64.NoSh12, sh12s[0], "Sh12.Shrink(LSL12) = %v", sh12s)
}

// TestBarrierDomainGen — domains from the eight spellings, shrink to Sy.
func TestBarrierDomainGen(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 200, BarrierDomain(rnd), func(d arm64.BarrierDomain) bool {
		return d >= arm64.Oshst && d <= arm64.Sy
	})

	domains := slices.Collect(BarrierDomain(arb.Rnd(1)).Shrink(arm64.Ish))
	require.Len(t, domains, 1, "BarrierDomain.Shrink(Ish) = %v", domains)
	require.Equal(t, arm64.Sy, domains[0], "BarrierDomain.Shrink(Ish) = %v", domains)
}

// TestSysregGen — every sample is a name of the arch table; the sorted
// order makes both the generation and the shrink target deterministic.
func TestSysregGen(t *testing.T) {
	names := sysregNamesSorted()
	require.NotEmpty(t, names, "the sysreg table is empty")
	require.True(t, slices.IsSorted(names), "sysregNamesSorted() is not sorted")

	valid := map[string]bool{}
	for _, name := range arm64.SysregNames() {
		valid[name] = true
	}

	rnd := arb.Rnd(42)
	ohsnap.Check(t, 300, Sysreg(rnd), func(name string) bool {
		return valid[name]
	})

	shrunk := slices.Collect(Sysreg(arb.Rnd(1)).Shrink(names[len(names)-1]))
	require.Len(t, shrunk, 1, "Sysreg.Shrink(last) = %v", shrunk)
	require.Equal(t, names[0], shrunk[0], "Sysreg.Shrink(last) = %v", shrunk)
}

// TestOffGen — aligned offsets, shrinking preserves alignment.
func TestOffGen(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 300, Off(rnd), func(o arm64.Off) bool {
		return o >= 0 && o <= 0x3ffc && o%4 == 0
	})
	for _, o := range offShrunk(0x180, 3) {
		require.True(t, o%8 == 0 && o >= 0, "offShrunk(0x180) → %d broke alignment", o)
	}
}

// --- instruction generators ------------------------------------------

// instrCase — one family: name + instruction generator.
type instrCase struct {
	name string
	gen  func() arm64.Instr
}

func newInstrCase(name string, gen func() arm64.Instr) instrCase {
	return instrCase{
		name: name,
		gen:  gen,
	}
}

// instrCases — all families.
func instrCases(rnd *rand.Rand) []instrCase {
	return []instrCase{
		newInstrCase("Ret", func() arm64.Instr {
			return ohsnap.First(Ret(rnd).Generate()).Instr()
		}),
		newInstrCase("Svc", func() arm64.Instr {
			return ohsnap.First(Svc(rnd).Generate()).Instr()
		}),
		newInstrCase("Brk", func() arm64.Instr {
			return ohsnap.First(Brk(rnd).Generate()).Instr()
		}),
		newInstrCase("Movz", func() arm64.Instr {
			return ohsnap.First(Movz(rnd).Generate()).Instr()
		}),
		newInstrCase("Movk", func() arm64.Instr {
			return ohsnap.First(Movk(rnd).Generate()).Instr()
		}),
		newInstrCase("AddImm", func() arm64.Instr {
			return ohsnap.First(AddImm(rnd).Generate()).Instr()
		}),
		newInstrCase("SubImm", func() arm64.Instr {
			return ohsnap.First(SubImm(rnd).Generate()).Instr()
		}),
		newInstrCase("AddShift", func() arm64.Instr {
			return ohsnap.First(AddShift(rnd).Generate()).Instr()
		}),
		newInstrCase("SubShift", func() arm64.Instr {
			return ohsnap.First(SubShift(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldr", func() arm64.Instr {
			return ohsnap.First(Ldr(rnd).Generate()).Instr()
		}),
		newInstrCase("Str", func() arm64.Instr {
			return ohsnap.First(Str(rnd).Generate()).Instr()
		}),
		newInstrCase("Nop", func() arm64.Instr {
			return ohsnap.First(Nop(rnd).Generate()).Instr()
		}),
		newInstrCase("Isb", func() arm64.Instr {
			return ohsnap.First(Isb(rnd).Generate()).Instr()
		}),
		newInstrCase("Dsb", func() arm64.Instr {
			return ohsnap.First(Dsb(rnd).Generate()).Instr()
		}),
		newInstrCase("Dmb", func() arm64.Instr {
			return ohsnap.First(Dmb(rnd).Generate()).Instr()
		}),
		newInstrCase("Smc", func() arm64.Instr {
			return ohsnap.First(Smc(rnd).Generate()).Instr()
		}),
		newInstrCase("Br", func() arm64.Instr {
			return ohsnap.First(Br(rnd).Generate()).Instr()
		}),
		newInstrCase("Blr", func() arm64.Instr {
			return ohsnap.First(Blr(rnd).Generate()).Instr()
		}),
		newInstrCase("Movn", func() arm64.Instr {
			return ohsnap.First(Movn(rnd).Generate()).Instr()
		}),
		newInstrCase("Adc", func() arm64.Instr {
			return ohsnap.First(Adc(rnd).Generate()).Instr()
		}),
		newInstrCase("Smulh", func() arm64.Instr {
			return ohsnap.First(Smulh(rnd).Generate()).Instr()
		}),
		newInstrCase("Umulh", func() arm64.Instr {
			return ohsnap.First(Umulh(rnd).Generate()).Instr()
		}),
		newInstrCase("Rev", func() arm64.Instr {
			return ohsnap.First(Rev(rnd).Generate()).Instr()
		}),
		newInstrCase("Rev16", func() arm64.Instr {
			return ohsnap.First(Rev16(rnd).Generate()).Instr()
		}),
		newInstrCase("Rev32", func() arm64.Instr {
			return ohsnap.First(Rev32(rnd).Generate()).Instr()
		}),
		newInstrCase("Cls", func() arm64.Instr {
			return ohsnap.First(Cls(rnd).Generate()).Instr()
		}),
		newInstrCase("Clz", func() arm64.Instr {
			return ohsnap.First(Clz(rnd).Generate()).Instr()
		}),
		newInstrCase("Rbit", func() arm64.Instr {
			return ohsnap.First(Rbit(rnd).Generate()).Instr()
		}),
		newInstrCase("Sdiv", func() arm64.Instr {
			return ohsnap.First(Sdiv(rnd).Generate()).Instr()
		}),
		newInstrCase("Udiv", func() arm64.Instr {
			return ohsnap.First(Udiv(rnd).Generate()).Instr()
		}),
		newInstrCase("LslReg", func() arm64.Instr {
			return ohsnap.First(LslReg(rnd).Generate()).Instr()
		}),
		newInstrCase("LsrReg", func() arm64.Instr {
			return ohsnap.First(LsrReg(rnd).Generate()).Instr()
		}),
		newInstrCase("AsrReg", func() arm64.Instr {
			return ohsnap.First(AsrReg(rnd).Generate()).Instr()
		}),
		newInstrCase("RorReg", func() arm64.Instr {
			return ohsnap.First(RorReg(rnd).Generate()).Instr()
		}),
		newInstrCase("Mrs", func() arm64.Instr {
			return ohsnap.First(Mrs(rnd).Generate()).Instr()
		}),
		newInstrCase("Msr", func() arm64.Instr {
			return ohsnap.First(Msr(rnd).Generate()).Instr()
		}),
		newInstrCase("B", func() arm64.Instr {
			return ohsnap.First(B(rnd).Generate()).Instr()
		}),
		newInstrCase("Bl", func() arm64.Instr {
			return ohsnap.First(Bl(rnd).Generate()).Instr()
		}),
		newInstrCase("Bcond", func() arm64.Instr {
			return ohsnap.First(Bcond(rnd).Generate()).Instr()
		}),
		newInstrCase("Cbz", func() arm64.Instr {
			return ohsnap.First(Cbz(rnd).Generate()).Instr()
		}),
		newInstrCase("Cbnz", func() arm64.Instr {
			return ohsnap.First(Cbnz(rnd).Generate()).Instr()
		}),
		newInstrCase("Tbz", func() arm64.Instr {
			return ohsnap.First(Tbz(rnd).Generate()).Instr()
		}),
		newInstrCase("Adr", func() arm64.Instr {
			return ohsnap.First(Adr(rnd).Generate()).Instr()
		}),
		newInstrCase("Adrp", func() arm64.Instr {
			return ohsnap.First(Adrp(rnd).Generate()).Instr()
		}),
		newInstrCase("AddsImm", func() arm64.Instr {
			return ohsnap.First(AddsImm(rnd).Generate()).Instr()
		}),
		newInstrCase("AddsShift", func() arm64.Instr {
			return ohsnap.First(AddsShift(rnd).Generate()).Instr()
		}),
		newInstrCase("AddExt", func() arm64.Instr {
			return ohsnap.First(AddExt(rnd).Generate()).Instr()
		}),
		newInstrCase("AddsExt", func() arm64.Instr {
			return ohsnap.First(AddsExt(rnd).Generate()).Instr()
		}),
		newInstrCase("SubExt", func() arm64.Instr {
			return ohsnap.First(SubExt(rnd).Generate()).Instr()
		}),
		newInstrCase("SubsExt", func() arm64.Instr {
			return ohsnap.First(SubsExt(rnd).Generate()).Instr()
		}),
		newInstrCase("AndShift", func() arm64.Instr {
			return ohsnap.First(AndShift(rnd).Generate()).Instr()
		}),
		newInstrCase("AndsShift", func() arm64.Instr {
			return ohsnap.First(AndsShift(rnd).Generate()).Instr()
		}),
		newInstrCase("OrrShift", func() arm64.Instr {
			return ohsnap.First(OrrShift(rnd).Generate()).Instr()
		}),
		newInstrCase("EorShift", func() arm64.Instr {
			return ohsnap.First(EorShift(rnd).Generate()).Instr()
		}),
		newInstrCase("OrnShift", func() arm64.Instr {
			return ohsnap.First(OrnShift(rnd).Generate()).Instr()
		}),
		newInstrCase("EonShift", func() arm64.Instr {
			return ohsnap.First(EonShift(rnd).Generate()).Instr()
		}),
		newInstrCase("BicShift", func() arm64.Instr {
			return ohsnap.First(BicShift(rnd).Generate()).Instr()
		}),
		newInstrCase("BicsShift", func() arm64.Instr {
			return ohsnap.First(BicsShift(rnd).Generate()).Instr()
		}),
		newInstrCase("AndImm", func() arm64.Instr {
			return ohsnap.First(AndImm(rnd).Generate()).Instr()
		}),
		newInstrCase("OrrImm", func() arm64.Instr {
			return ohsnap.First(OrrImm(rnd).Generate()).Instr()
		}),
		newInstrCase("EorImm", func() arm64.Instr {
			return ohsnap.First(EorImm(rnd).Generate()).Instr()
		}),
		newInstrCase("AndsImm", func() arm64.Instr {
			return ohsnap.First(AndsImm(rnd).Generate()).Instr()
		}),
		newInstrCase("Bfm", func() arm64.Instr {
			return ohsnap.First(Bfm(rnd).Generate()).Instr()
		}),
		newInstrCase("Sbfm", func() arm64.Instr {
			return ohsnap.First(Sbfm(rnd).Generate()).Instr()
		}),
		newInstrCase("Ubfm", func() arm64.Instr {
			return ohsnap.First(Ubfm(rnd).Generate()).Instr()
		}),
		newInstrCase("Extr", func() arm64.Instr {
			return ohsnap.First(Extr(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldrb", func() arm64.Instr {
			return ohsnap.First(Ldrb(rnd).Generate()).Instr()
		}),
		newInstrCase("Strb", func() arm64.Instr {
			return ohsnap.First(Strb(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldrh", func() arm64.Instr {
			return ohsnap.First(Ldrh(rnd).Generate()).Instr()
		}),
		newInstrCase("Strh", func() arm64.Instr {
			return ohsnap.First(Strh(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldrsb", func() arm64.Instr {
			return ohsnap.First(Ldrsb(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldrsh", func() arm64.Instr {
			return ohsnap.First(Ldrsh(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldrsw", func() arm64.Instr {
			return ohsnap.First(Ldrsw(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldur", func() arm64.Instr {
			return ohsnap.First(Ldur(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldurb", func() arm64.Instr {
			return ohsnap.First(Ldurb(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldurh", func() arm64.Instr {
			return ohsnap.First(Ldurh(rnd).Generate()).Instr()
		}),
		newInstrCase("Stur", func() arm64.Instr {
			return ohsnap.First(Stur(rnd).Generate()).Instr()
		}),
		newInstrCase("Sturb", func() arm64.Instr {
			return ohsnap.First(Sturb(rnd).Generate()).Instr()
		}),
		newInstrCase("Sturh", func() arm64.Instr {
			return ohsnap.First(Sturh(rnd).Generate()).Instr()
		}),
		newInstrCase("LdrF", func() arm64.Instr {
			return ohsnap.First(LdrF(rnd).Generate()).Instr()
		}),
		newInstrCase("StrF", func() arm64.Instr {
			return ohsnap.First(StrF(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldp", func() arm64.Instr {
			return ohsnap.First(Ldp(rnd).Generate()).Instr()
		}),
		newInstrCase("Stp", func() arm64.Instr {
			return ohsnap.First(Stp(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldpsw", func() arm64.Instr {
			return ohsnap.First(Ldpsw(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldar", func() arm64.Instr {
			return ohsnap.First(Ldar(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldaxr", func() arm64.Instr {
			return ohsnap.First(Ldaxr(rnd).Generate()).Instr()
		}),
		newInstrCase("Stlr", func() arm64.Instr {
			return ohsnap.First(Stlr(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldarb", func() arm64.Instr {
			return ohsnap.First(Ldarb(rnd).Generate()).Instr()
		}),
		newInstrCase("Ldaxrb", func() arm64.Instr {
			return ohsnap.First(Ldaxrb(rnd).Generate()).Instr()
		}),
		newInstrCase("Stlrb", func() arm64.Instr {
			return ohsnap.First(Stlrb(rnd).Generate()).Instr()
		}),
		newInstrCase("Stlxr", func() arm64.Instr {
			return ohsnap.First(Stlxr(rnd).Generate()).Instr()
		}),
		newInstrCase("Stxrb", func() arm64.Instr {
			return ohsnap.First(Stxrb(rnd).Generate()).Instr()
		}),
		newInstrCase("Stlxrb", func() arm64.Instr {
			return ohsnap.First(Stlxrb(rnd).Generate()).Instr()
		}),
		newInstrCase("Prfm", func() arm64.Instr {
			return ohsnap.First(Prfm(rnd).Generate()).Instr()
		}),
	}
}

// TestInstrGenValid — a generated instruction always encodes
// (the generator produces no invalid combinations — otherwise a constructor error).
func TestInstrGenValid(t *testing.T) {
	rnd := arb.Rnd(42)
	for _, c := range instrCases(rnd) {
		var buf bytes.Buffer
		for range 300 {
			buf.Reset()
			in := c.gen()
			_, err := in.Encode(&buf)
			require.NoError(t, err, "%s: Encode (%s)", c.name, in.ObjDump(disasm.DefaultViewCtx()))
		}
	}
}

// TestInstrGenText — String() of the parameters = ObjDump of the instruction.
func TestInstrGenText(t *testing.T) {
	rnd := arb.Rnd(42)
	require.NotEmpty(t, ohsnap.First(Movz(rnd).Generate()).String(), "MovzParams.String() is empty")
	got := ohsnap.First(Ldr(rnd).Generate()).String()
	require.True(t, strings.HasPrefix(got, "ldr "), "LdrParams.String() = %q", got)
}

// TestInstrGenShrinkValid — every shrink candidate is a valid instruction
// (component-wise shrinking does not break contextual invariants).
func TestInstrGenShrinkValid(t *testing.T) {
	rnd := arb.Rnd(9)
	// Movz/Movk/Movn: shrinking a W-form never leaves hw>=2.
	movz := Movz(rnd)
	for range 200 {
		p := ohsnap.First(movz.Generate())
		for s := range movz.Shrink(p) {
			require.True(
				t,
				s.Rd.Is64() || s.Hw <= arm64.Hw1,
				"shrink broke the invariant: %s",
				s,
			)
			s.Instr() // the constructor error path is unreachable by construction of the shrink
		}
	}

	movn := Movn(rnd)
	for range 200 {
		p := ohsnap.First(movn.Generate())
		for s := range movn.Shrink(p) {
			require.True(
				t,
				s.Rd.Is64() || s.Hw <= arm64.Hw1,
				"shrink broke the invariant: %s",
				s,
			)
			s.Instr()
		}
	}

	// AddImm: the width of rd/rn matches after shrinking.
	add := AddImm(rnd)
	for range 200 {
		p := ohsnap.First(add.Generate())
		for s := range add.Shrink(p) {
			require.Equal(t, s.Rd.Is64(), s.Rn.Is64(), "shrink broke the width: %s", s)
			s.Instr()
		}
	}

	// Sdiv: the width of rd/rn/rm matches after shrinking.
	sdiv := Sdiv(rnd)
	for range 200 {
		p := ohsnap.First(sdiv.Generate())
		for s := range sdiv.Shrink(p) {
			require.Equal(
				t,
				s.Rd.Is64(),
				s.Rn.Is64() && s.Rm.Is64(),
				"shrink broke the width: %s",
				s,
			)
			s.Instr()
		}
	}

	// Ldr: the offset is aligned to the width of rt after shrinking.
	ldr := Ldr(rnd)
	for range 200 {
		p := ohsnap.First(ldr.Generate())
		for s := range ldr.Shrink(p) {
			align := int64(4)
			if s.Rt.Is64() {
				align = 8
			}

			require.Zero(t, int64(s.Off)%align, "shrink broke the alignment: %s", s)
			s.Instr()
		}
	}
}
