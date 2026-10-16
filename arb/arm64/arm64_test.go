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
		seen[Reg(rnd).Generate().String()] = true
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
	for s := range i12.Shrink(i12.Generate()) {
		n := mustImmValue(t, s)
		require.True(t, n >= 0 && n <= 0xfff, "Imm12 shrink out of range: %v", s)
	}

	for s := range i16.Shrink(i16.Generate()) {
		n := mustImmValue(t, s)
		require.True(t, n >= 0 && n <= 0xffff, "Imm16 shrink out of range: %v", s)
	}

	n := mustImmValue(t, Imm6(rnd).Generate())
	require.GreaterOrEqual(t, n, int64(0), "Imm6: immValue failed to parse String()")
}

func TestImmGenBounds(t *testing.T) {
	// With a fixed seed the bounds are reachable.
	rnd := arb.Rnd(7)
	min12, max12 := int64(1<<62), int64(-1)
	for range 2000 {
		n := mustImmValue(t, Imm12(rnd).Generate())
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
	require.NotEmpty(t, Movz(rnd).Generate().String(), "MovzParams.String() is empty")
	got := Ldr(rnd).Generate().String()
	require.True(t, strings.HasPrefix(got, "ldr "), "LdrParams.String() = %q", got)
}

// TestInstrGenShrinkValid — every shrink candidate is a valid instruction
// (component-wise shrinking does not break contextual invariants).
func TestInstrGenShrinkValid(t *testing.T) {
	rnd := arb.Rnd(9)
	// Movz/Movk/Movn: shrinking a W-form never leaves hw>=2.
	movz := Movz(rnd)
	for range 200 {
		p := movz.Generate()
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
		p := movn.Generate()
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
		p := add.Generate()
		for s := range add.Shrink(p) {
			require.Equal(t, s.Rd.Is64(), s.Rn.Is64(), "shrink broke the width: %s", s)
			s.Instr()
		}
	}

	// Sdiv: the width of rd/rn/rm matches after shrinking.
	sdiv := Sdiv(rnd)
	for range 200 {
		p := sdiv.Generate()
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
		p := ldr.Generate()
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
			return Ret(rnd).Generate().Instr()
		}),
		newInstrCase("Svc", func() arm64.Instr {
			return Svc(rnd).Generate().Instr()
		}),
		newInstrCase("Brk", func() arm64.Instr {
			return Brk(rnd).Generate().Instr()
		}),
		newInstrCase("Movz", func() arm64.Instr {
			return Movz(rnd).Generate().Instr()
		}),
		newInstrCase("Movk", func() arm64.Instr {
			return Movk(rnd).Generate().Instr()
		}),
		newInstrCase("AddImm", func() arm64.Instr {
			return AddImm(rnd).Generate().Instr()
		}),
		newInstrCase("SubImm", func() arm64.Instr {
			return SubImm(rnd).Generate().Instr()
		}),
		newInstrCase("AddShift", func() arm64.Instr {
			return AddShift(rnd).Generate().Instr()
		}),
		newInstrCase("SubShift", func() arm64.Instr {
			return SubShift(rnd).Generate().Instr()
		}),
		newInstrCase("Ldr", func() arm64.Instr {
			return Ldr(rnd).Generate().Instr()
		}),
		newInstrCase("Str", func() arm64.Instr {
			return Str(rnd).Generate().Instr()
		}),
		newInstrCase("Nop", func() arm64.Instr {
			return Nop(rnd).Generate().Instr()
		}),
		newInstrCase("Isb", func() arm64.Instr {
			return Isb(rnd).Generate().Instr()
		}),
		newInstrCase("Dsb", func() arm64.Instr {
			return Dsb(rnd).Generate().Instr()
		}),
		newInstrCase("Dmb", func() arm64.Instr {
			return Dmb(rnd).Generate().Instr()
		}),
		newInstrCase("Smc", func() arm64.Instr {
			return Smc(rnd).Generate().Instr()
		}),
		newInstrCase("Br", func() arm64.Instr {
			return Br(rnd).Generate().Instr()
		}),
		newInstrCase("Blr", func() arm64.Instr {
			return Blr(rnd).Generate().Instr()
		}),
		newInstrCase("Movn", func() arm64.Instr {
			return Movn(rnd).Generate().Instr()
		}),
		newInstrCase("Adc", func() arm64.Instr {
			return Adc(rnd).Generate().Instr()
		}),
		newInstrCase("Smulh", func() arm64.Instr {
			return Smulh(rnd).Generate().Instr()
		}),
		newInstrCase("Umulh", func() arm64.Instr {
			return Umulh(rnd).Generate().Instr()
		}),
		newInstrCase("Rev", func() arm64.Instr {
			return Rev(rnd).Generate().Instr()
		}),
		newInstrCase("Rev16", func() arm64.Instr {
			return Rev16(rnd).Generate().Instr()
		}),
		newInstrCase("Rev32", func() arm64.Instr {
			return Rev32(rnd).Generate().Instr()
		}),
		newInstrCase("Cls", func() arm64.Instr {
			return Cls(rnd).Generate().Instr()
		}),
		newInstrCase("Clz", func() arm64.Instr {
			return Clz(rnd).Generate().Instr()
		}),
		newInstrCase("Rbit", func() arm64.Instr {
			return Rbit(rnd).Generate().Instr()
		}),
		newInstrCase("Sdiv", func() arm64.Instr {
			return Sdiv(rnd).Generate().Instr()
		}),
		newInstrCase("Udiv", func() arm64.Instr {
			return Udiv(rnd).Generate().Instr()
		}),
		newInstrCase("LslReg", func() arm64.Instr {
			return LslReg(rnd).Generate().Instr()
		}),
		newInstrCase("LsrReg", func() arm64.Instr {
			return LsrReg(rnd).Generate().Instr()
		}),
		newInstrCase("AsrReg", func() arm64.Instr {
			return AsrReg(rnd).Generate().Instr()
		}),
		newInstrCase("RorReg", func() arm64.Instr {
			return RorReg(rnd).Generate().Instr()
		}),
		newInstrCase("Mrs", func() arm64.Instr {
			return Mrs(rnd).Generate().Instr()
		}),
		newInstrCase("Msr", func() arm64.Instr {
			return Msr(rnd).Generate().Instr()
		}),
		newInstrCase("B", func() arm64.Instr {
			return B(rnd).Generate().Instr()
		}),
		newInstrCase("Bl", func() arm64.Instr {
			return Bl(rnd).Generate().Instr()
		}),
		newInstrCase("Bcond", func() arm64.Instr {
			return Bcond(rnd).Generate().Instr()
		}),
		newInstrCase("Cbz", func() arm64.Instr {
			return Cbz(rnd).Generate().Instr()
		}),
		newInstrCase("Cbnz", func() arm64.Instr {
			return Cbnz(rnd).Generate().Instr()
		}),
		newInstrCase("Tbz", func() arm64.Instr {
			return Tbz(rnd).Generate().Instr()
		}),
		newInstrCase("Adr", func() arm64.Instr {
			return Adr(rnd).Generate().Instr()
		}),
		newInstrCase("Adrp", func() arm64.Instr {
			return Adrp(rnd).Generate().Instr()
		}),
		newInstrCase("AddsImm", func() arm64.Instr {
			return AddsImm(rnd).Generate().Instr()
		}),
		newInstrCase("AddsShift", func() arm64.Instr {
			return AddsShift(rnd).Generate().Instr()
		}),
		newInstrCase("AddExt", func() arm64.Instr {
			return AddExt(rnd).Generate().Instr()
		}),
		newInstrCase("AddsExt", func() arm64.Instr {
			return AddsExt(rnd).Generate().Instr()
		}),
		newInstrCase("SubExt", func() arm64.Instr {
			return SubExt(rnd).Generate().Instr()
		}),
		newInstrCase("SubsExt", func() arm64.Instr {
			return SubsExt(rnd).Generate().Instr()
		}),
		newInstrCase("AndShift", func() arm64.Instr {
			return AndShift(rnd).Generate().Instr()
		}),
		newInstrCase("AndsShift", func() arm64.Instr {
			return AndsShift(rnd).Generate().Instr()
		}),
		newInstrCase("OrrShift", func() arm64.Instr {
			return OrrShift(rnd).Generate().Instr()
		}),
		newInstrCase("EorShift", func() arm64.Instr {
			return EorShift(rnd).Generate().Instr()
		}),
		newInstrCase("OrnShift", func() arm64.Instr {
			return OrnShift(rnd).Generate().Instr()
		}),
		newInstrCase("EonShift", func() arm64.Instr {
			return EonShift(rnd).Generate().Instr()
		}),
		newInstrCase("BicShift", func() arm64.Instr {
			return BicShift(rnd).Generate().Instr()
		}),
		newInstrCase("BicsShift", func() arm64.Instr {
			return BicsShift(rnd).Generate().Instr()
		}),
		newInstrCase("AndImm", func() arm64.Instr {
			return AndImm(rnd).Generate().Instr()
		}),
		newInstrCase("OrrImm", func() arm64.Instr {
			return OrrImm(rnd).Generate().Instr()
		}),
		newInstrCase("EorImm", func() arm64.Instr {
			return EorImm(rnd).Generate().Instr()
		}),
		newInstrCase("AndsImm", func() arm64.Instr {
			return AndsImm(rnd).Generate().Instr()
		}),
		newInstrCase("Bfm", func() arm64.Instr {
			return Bfm(rnd).Generate().Instr()
		}),
		newInstrCase("Sbfm", func() arm64.Instr {
			return Sbfm(rnd).Generate().Instr()
		}),
		newInstrCase("Ubfm", func() arm64.Instr {
			return Ubfm(rnd).Generate().Instr()
		}),
		newInstrCase("Extr", func() arm64.Instr {
			return Extr(rnd).Generate().Instr()
		}),
		newInstrCase("Ldrb", func() arm64.Instr {
			return Ldrb(rnd).Generate().Instr()
		}),
		newInstrCase("Strb", func() arm64.Instr {
			return Strb(rnd).Generate().Instr()
		}),
		newInstrCase("Ldrh", func() arm64.Instr {
			return Ldrh(rnd).Generate().Instr()
		}),
		newInstrCase("Strh", func() arm64.Instr {
			return Strh(rnd).Generate().Instr()
		}),
		newInstrCase("Ldrsb", func() arm64.Instr {
			return Ldrsb(rnd).Generate().Instr()
		}),
		newInstrCase("Ldrsh", func() arm64.Instr {
			return Ldrsh(rnd).Generate().Instr()
		}),
		newInstrCase("Ldrsw", func() arm64.Instr {
			return Ldrsw(rnd).Generate().Instr()
		}),
		newInstrCase("Ldur", func() arm64.Instr {
			return Ldur(rnd).Generate().Instr()
		}),
		newInstrCase("Ldurb", func() arm64.Instr {
			return Ldurb(rnd).Generate().Instr()
		}),
		newInstrCase("Ldurh", func() arm64.Instr {
			return Ldurh(rnd).Generate().Instr()
		}),
		newInstrCase("Stur", func() arm64.Instr {
			return Stur(rnd).Generate().Instr()
		}),
		newInstrCase("Sturb", func() arm64.Instr {
			return Sturb(rnd).Generate().Instr()
		}),
		newInstrCase("Sturh", func() arm64.Instr {
			return Sturh(rnd).Generate().Instr()
		}),
		newInstrCase("LdrF", func() arm64.Instr {
			return LdrF(rnd).Generate().Instr()
		}),
		newInstrCase("StrF", func() arm64.Instr {
			return StrF(rnd).Generate().Instr()
		}),
		newInstrCase("Ldp", func() arm64.Instr {
			return Ldp(rnd).Generate().Instr()
		}),
		newInstrCase("Stp", func() arm64.Instr {
			return Stp(rnd).Generate().Instr()
		}),
		newInstrCase("Ldpsw", func() arm64.Instr {
			return Ldpsw(rnd).Generate().Instr()
		}),
		newInstrCase("Ldar", func() arm64.Instr {
			return Ldar(rnd).Generate().Instr()
		}),
		newInstrCase("Ldaxr", func() arm64.Instr {
			return Ldaxr(rnd).Generate().Instr()
		}),
		newInstrCase("Stlr", func() arm64.Instr {
			return Stlr(rnd).Generate().Instr()
		}),
		newInstrCase("Ldarb", func() arm64.Instr {
			return Ldarb(rnd).Generate().Instr()
		}),
		newInstrCase("Ldaxrb", func() arm64.Instr {
			return Ldaxrb(rnd).Generate().Instr()
		}),
		newInstrCase("Stlrb", func() arm64.Instr {
			return Stlrb(rnd).Generate().Instr()
		}),
		newInstrCase("Stlxr", func() arm64.Instr {
			return Stlxr(rnd).Generate().Instr()
		}),
		newInstrCase("Stxrb", func() arm64.Instr {
			return Stxrb(rnd).Generate().Instr()
		}),
		newInstrCase("Stlxrb", func() arm64.Instr {
			return Stlxrb(rnd).Generate().Instr()
		}),
		newInstrCase("Prfm", func() arm64.Instr {
			return Prfm(rnd).Generate().Instr()
		}),
		newInstrCase("Fadd", func() arm64.Instr {
			return Fadd(rnd).Generate().Instr()
		}),
		newInstrCase("Fsub", func() arm64.Instr {
			return Fsub(rnd).Generate().Instr()
		}),
		newInstrCase("Fmul", func() arm64.Instr {
			return Fmul(rnd).Generate().Instr()
		}),
		newInstrCase("Fdiv", func() arm64.Instr {
			return Fdiv(rnd).Generate().Instr()
		}),
		newInstrCase("Fmax", func() arm64.Instr {
			return Fmax(rnd).Generate().Instr()
		}),
		newInstrCase("Fmin", func() arm64.Instr {
			return Fmin(rnd).Generate().Instr()
		}),
		newInstrCase("Fcmp", func() arm64.Instr {
			return Fcmp(rnd).Generate().Instr()
		}),
		newInstrCase("FcmpZero", func() arm64.Instr {
			return FcmpZero(rnd).Generate().Instr()
		}),
		newInstrCase("Fneg", func() arm64.Instr {
			return Fneg(rnd).Generate().Instr()
		}),
		newInstrCase("Fmov", func() arm64.Instr {
			return Fmov(rnd).Generate().Instr()
		}),
		newInstrCase("Fcvt", func() arm64.Instr {
			return Fcvt(rnd).Generate().Instr()
		}),
		newInstrCase("Fmadd", func() arm64.Instr {
			return Fmadd(rnd).Generate().Instr()
		}),
		newInstrCase("Fnmsub", func() arm64.Instr {
			return Fnmsub(rnd).Generate().Instr()
		}),
		newInstrCase("FmovFromGpr", func() arm64.Instr {
			return FmovFromGpr(rnd).Generate().Instr()
		}),
		newInstrCase("FmovToGpr", func() arm64.Instr {
			return FmovToGpr(rnd).Generate().Instr()
		}),
		newInstrCase("Fcvtzs", func() arm64.Instr {
			return Fcvtzs(rnd).Generate().Instr()
		}),
		newInstrCase("Fcvtzu", func() arm64.Instr {
			return Fcvtzu(rnd).Generate().Instr()
		}),
		newInstrCase("Scvtf", func() arm64.Instr {
			return Scvtf(rnd).Generate().Instr()
		}),
		newInstrCase("Ucvtf", func() arm64.Instr {
			return Ucvtf(rnd).Generate().Instr()
		}),
		newInstrCase("FmovImm", func() arm64.Instr {
			return FmovImm(rnd).Generate().Instr()
		}),
		newInstrCase("Add", func() arm64.Instr {
			return Add(rnd).Generate().Instr()
		}),
		newInstrCase("Addp", func() arm64.Instr {
			return Addp(rnd).Generate().Instr()
		}),
		newInstrCase("And", func() arm64.Instr {
			return And(rnd).Generate().Instr()
		}),
		newInstrCase("Bic", func() arm64.Instr {
			return Bic(rnd).Generate().Instr()
		}),
		newInstrCase("Bif", func() arm64.Instr {
			return Bif(rnd).Generate().Instr()
		}),
		newInstrCase("Bit", func() arm64.Instr {
			return Bit(rnd).Generate().Instr()
		}),
		newInstrCase("Bsl", func() arm64.Instr {
			return Bsl(rnd).Generate().Instr()
		}),
		newInstrCase("Cmeq", func() arm64.Instr {
			return Cmeq(rnd).Generate().Instr()
		}),
		newInstrCase("Cmge", func() arm64.Instr {
			return Cmge(rnd).Generate().Instr()
		}),
		newInstrCase("Cmtst", func() arm64.Instr {
			return Cmtst(rnd).Generate().Instr()
		}),
		newInstrCase("Eor", func() arm64.Instr {
			return Eor(rnd).Generate().Instr()
		}),
		newInstrCase("Orn", func() arm64.Instr {
			return Orn(rnd).Generate().Instr()
		}),
		newInstrCase("Orr", func() arm64.Instr {
			return Orr(rnd).Generate().Instr()
		}),
		newInstrCase("Saddw", func() arm64.Instr {
			return Saddw(rnd).Generate().Instr()
		}),
		newInstrCase("Ssubw", func() arm64.Instr {
			return Ssubw(rnd).Generate().Instr()
		}),
		newInstrCase("Uaddw", func() arm64.Instr {
			return Uaddw(rnd).Generate().Instr()
		}),
		newInstrCase("Usubw", func() arm64.Instr {
			return Usubw(rnd).Generate().Instr()
		}),
		newInstrCase("Sqrshl", func() arm64.Instr {
			return Sqrshl(rnd).Generate().Instr()
		}),
		newInstrCase("Abs", func() arm64.Instr {
			return Abs(rnd).Generate().Instr()
		}),
		newInstrCase("Cnt", func() arm64.Instr {
			return Cnt(rnd).Generate().Instr()
		}),
		newInstrCase("Not", func() arm64.Instr {
			return Not(rnd).Generate().Instr()
		}),
		newInstrCase("Rev32V", func() arm64.Instr {
			return Rev32V(rnd).Generate().Instr()
		}),
		newInstrCase("RbitV", func() arm64.Instr {
			return RbitV(rnd).Generate().Instr()
		}),
		newInstrCase("Uaddlv", func() arm64.Instr {
			return Uaddlv(rnd).Generate().Instr()
		}),
		newInstrCase("MovSimd", func() arm64.Instr {
			return MovSimd(rnd).Generate().Instr()
		}),
		newInstrCase("DupGen", func() arm64.Instr {
			return DupGen(rnd).Generate().Instr()
		}),
		newInstrCase("Shl", func() arm64.Instr {
			return Shl(rnd).Generate().Instr()
		}),
		newInstrCase("Sri", func() arm64.Instr {
			return Sri(rnd).Generate().Instr()
		}),
		newInstrCase("Sshr", func() arm64.Instr {
			return Sshr(rnd).Generate().Instr()
		}),
		newInstrCase("Ushr", func() arm64.Instr {
			return Ushr(rnd).Generate().Instr()
		}),
		newInstrCase("Tbl", func() arm64.Instr {
			return Tbl(rnd).Generate().Instr()
		}),
		newInstrCase("Aese", func() arm64.Instr {
			return Aese(rnd).Generate().Instr()
		}),
		newInstrCase("Aesmc", func() arm64.Instr {
			return Aesmc(rnd).Generate().Instr()
		}),
	}
}
