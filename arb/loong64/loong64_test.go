package loong64

import (
	"bytes"
	"math/rand/v2"
	"slices"
	"testing"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/arb"
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/disasm"
)

// --- operand generators --------------------------------------------------

// TestRegGenProperty — samples in $r0..$r31, shrink to zero.
func TestRegGenProperty(t *testing.T) {
	rnd := arb.Rnd(42)
	ohsnap.Check(t, 500, Reg(rnd), func(r arch.Reg) bool {
		return r.Num() <= 31
	})
}

func TestRegGenDistribution(t *testing.T) {
	rnd := arb.Rnd(42)
	seen := map[string]bool{}
	for range 1000 {
		seen[ohsnap.First(Reg(rnd).Generate()).String()] = true
	}

	for _, want := range []string{"$zero", "$ra", "$tp", "$sp", "$a0", "$t0", "$r21", "$s8"} {
		require.True(t, seen[want], "%s was not seen in 1000 samples", want)
	}
}

func TestRegGenShrink(t *testing.T) {
	r9, err := arch.R(9)
	require.NoError(t, err)
	got := regShrunk(r9)
	require.Len(t, got, 2, "regShrunk($r9) = %v", got)
	require.Equal(t, arch.Zero, got[0])
	r4, err := arch.R(4)
	require.NoError(t, err)
	require.Equal(t, r4, got[1])
	got = regShrunk(arch.Zero)
	require.Empty(t, got, "regShrunk(zero) — shrinker self-loop")
	r1, err := arch.R(1)
	require.NoError(t, err)
	got = regShrunk(r1)
	require.Len(t, got, 1, "regShrunk($r1) = %v", got)
	require.Equal(t, arch.Zero, got[0])

	// the same through the Arbitrary interface
	require.Equal(t, regShrunk(r9), slices.Collect(Reg(arb.Rnd(7)).Shrink(r9)))
}

// TestImmGenProperty — every role: samples in range (the aligned roles
// in step), shrink candidates in range. One t.Run per role: each role
// type has its own checked constructor and range.
func TestImmGenProperty(t *testing.T) {
	t.Run("Imm12", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, Imm12(rnd), func(v arch.Imm12) bool {
			return v.Val() >= -2048 && v.Val() <= 2047
		})

		lo, err := arch.NewImm12(-2048)
		require.NoError(t, err)
		for s := range Imm12(rnd).Shrink(lo) {
			require.GreaterOrEqual(t, s.Val(), int64(-2048), "shrink out of range: %v", s)
			require.LessOrEqual(t, s.Val(), int64(2047), "shrink out of range: %v", s)
		}
	})

	t.Run("UImm12", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, UImm12(rnd), func(v arch.UImm12) bool {
			return v.Val() >= 0 && v.Val() <= 4095
		})
	})

	t.Run("Imm14", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, Imm14(rnd), func(v arch.Imm14) bool {
			return v.Val()%4 == 0 && v.Val() >= -16380 && v.Val() <= 16380
		})

		// the halved candidates that lost alignment are skipped: the
		// surviving candidates are aligned and in range
		hi, err := arch.NewImm14(16380)
		require.NoError(t, err)
		cs := slices.Collect(Imm14(rnd).Shrink(hi))
		require.NotEmpty(t, cs)
		for _, s := range cs {
			require.Equal(t, int64(0), s.Val()%4, "unaligned shrink candidate: %v", s)
			require.GreaterOrEqual(t, s.Val(), int64(-16380), "shrink out of range: %v", s)
		}
	})

	t.Run("Imm16", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, Imm16(rnd), func(v arch.Imm16) bool {
			return v.Val() >= -32768 && v.Val() <= 32767
		})
	})

	t.Run("Off16", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, Off16(rnd), func(v arch.Off16) bool {
			return v.Val()%4 == 0 && v.Val() >= -131068 && v.Val() <= 131068
		})
	})

	t.Run("Imm20", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, Imm20(rnd), func(v arch.Imm20) bool {
			return v.Val() >= -524288 && v.Val() <= 524287
		})
	})

	t.Run("UImm5", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, UImm5(rnd), func(v arch.UImm5) bool {
			return v.Val() >= 0 && v.Val() <= 31
		})
	})

	t.Run("UImm6", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, UImm6(rnd), func(v arch.UImm6) bool {
			return v.Val() >= 0 && v.Val() <= 63
		})
	})

	t.Run("UImm2", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 100, UImm2(rnd), func(v arch.UImm2) bool {
			return v.Val() >= 0 && v.Val() <= 3
		})
	})

	t.Run("UImm3", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 100, UImm3(rnd), func(v arch.UImm3) bool {
			return v.Val() >= 0 && v.Val() <= 7
		})
	})

	t.Run("Shift3", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 100, Shift3(rnd), func(v arch.Shift3) bool {
			return v.Val() >= 1 && v.Val() <= 4
		})
	})

	t.Run("UImm8", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, UImm8(rnd), func(v arch.UImm8) bool {
			return v.Val() >= 0 && v.Val() <= 255
		})
	})

	t.Run("UImm14", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, UImm14(rnd), func(v arch.UImm14) bool {
			return v.Val() >= 0 && v.Val() <= 16383
		})
	})

	t.Run("Code15", func(t *testing.T) {
		rnd := arb.Rnd(42)
		ohsnap.Check(t, 300, Code15(rnd), func(v arch.Code15) bool {
			return v.Val() >= 0 && v.Val() <= 32767
		})
	})
}

// --- instruction families ------------------------------------------------

// laInstrParam — parameters of any family (the generic bridge of
// newInstrCase: families have different parameter types, the sampler
// closes over the instantiation).
type laInstrParam interface {
	Instr() arch.Instr
	String() string
}

// instrCase — one family: a sampler (instruction + params text) and a
// shrink round (the candidates of one generated value as instructions).
type instrCase struct {
	name   string
	sample func() (arch.Instr, string)
	shrink func() []arch.Instr
}

// newInstrCase — a family entry: closes over the generic instantiation
// (Arbitrary is invariant — a common Arbitrary[laInstrParam] cannot be
// assembled, see newRvFamily in property_rv_test.go).
func newInstrCase[P laInstrParam](name string, mk func() ohsnap.Arbitrary[P]) instrCase {
	return instrCase{
		name: name,
		sample: func() (arch.Instr, string) {
			p := ohsnap.First(mk().Generate())
			return p.Instr(), p.String()
		},
		shrink: func() []arch.Instr {
			g := mk()
			cs := slices.Collect(g.Shrink(ohsnap.First(g.Generate())))
			out := make([]arch.Instr, len(cs))
			for i, c := range cs {
				out[i] = c.Instr()
			}

			return out
		},
	}
}

func instrCases(rnd *rand.Rand) []instrCase {
	return []instrCase{
		newInstrCase("Alu3R", func() ohsnap.Arbitrary[R3Params] { return Alu3R(rnd) }),
		newInstrCase("Alu2R", func() ohsnap.Arbitrary[R2Params] { return Alu2R(rnd) }),
		newInstrCase("AluImm12", func() ohsnap.Arbitrary[R2RoleParams[arch.Imm12]] {
			return AluImm12(rnd)
		}),
		newInstrCase("AluUImm12", func() ohsnap.Arbitrary[R2RoleParams[arch.UImm12]] {
			return AluUImm12(rnd)
		}),
		newInstrCase("AluImm16", func() ohsnap.Arbitrary[R2RoleParams[arch.Imm16]] {
			return AluImm16(rnd)
		}),
		newInstrCase("Imm20", func() ohsnap.Arbitrary[R1RoleParams[arch.Imm20]] {
			return Imm20Instr(rnd)
		}),
		newInstrCase("Code15", func() ohsnap.Arbitrary[CodeParams] { return Code15Instr(rnd) }),
		newInstrCase("Branch2", func() ohsnap.Arbitrary[Branch2Params] { return Branch2(rnd) }),
		newInstrCase("Branch1", func() ohsnap.Arbitrary[Branch1Params] { return Branch1(rnd) }),
		newInstrCase("Jump", func() ohsnap.Arbitrary[JumpParams] { return Jump(rnd) }),
		newInstrCase("Jirl", func() ohsnap.Arbitrary[R2RoleParams[arch.Off16]] {
			return Jirl(rnd)
		}),
		newInstrCase("LdSt", func() ohsnap.Arbitrary[R2RoleParams[arch.Imm12]] {
			return LdSt(rnd)
		}),
		newInstrCase("Ldptr", func() ohsnap.Arbitrary[R2RoleParams[arch.Imm14]] {
			return Ldptr(rnd)
		}),
		newInstrCase("LdxStx", func() ohsnap.Arbitrary[R3Params] { return LdxStx(rnd) }),
		newInstrCase("LdAcq", func() ohsnap.Arbitrary[R2Params] { return LdAcq(rnd) }),
		newInstrCase("Hints", func() ohsnap.Arbitrary[U5RI12Params] { return Hints(rnd) }),
		newInstrCase("Preldx", func() ohsnap.Arbitrary[U5RRParams] { return Preldx(rnd) }),
		newInstrCase("ShiftW", func() ohsnap.Arbitrary[R2RoleParams[arch.UImm5]] {
			return ShiftW(rnd)
		}),
		newInstrCase("ShiftD", func() ohsnap.Arbitrary[R2RoleParams[arch.UImm6]] {
			return ShiftD(rnd)
		}),
		newInstrCase("FieldW", func() ohsnap.Arbitrary[FieldWParams] { return FieldW(rnd) }),
		newInstrCase("FieldD", func() ohsnap.Arbitrary[FieldDParams] { return FieldD(rnd) }),
		newInstrCase("Alsl", func() ohsnap.Arbitrary[R3RoleParams[arch.Shift3]] {
			return Alsl(rnd)
		}),
		newInstrCase("BytepickW", func() ohsnap.Arbitrary[R3RoleParams[arch.UImm2]] {
			return BytepickW(rnd)
		}),
		newInstrCase("BytepickD", func() ohsnap.Arbitrary[R3RoleParams[arch.UImm3]] {
			return BytepickD(rnd)
		}),
		newInstrCase("Atomics", func() ohsnap.Arbitrary[R3Params] { return Atomics(rnd) }),
		newInstrCase("CsrRW", func() ohsnap.Arbitrary[R1RoleParams[arch.UImm14]] {
			return CsrRW(rnd)
		}),
		newInstrCase("CsrXchg", func() ohsnap.Arbitrary[R2RoleParams[arch.UImm14]] {
			return CsrXchg(rnd)
		}),
		newInstrCase("IoCsr", func() ohsnap.Arbitrary[R2Params] { return IoCsr(rnd) }),
		newInstrCase("Lddir", func() ohsnap.Arbitrary[R2RoleParams[arch.UImm8]] {
			return Lddir(rnd)
		}),
		newInstrCase("Ldpte", func() ohsnap.Arbitrary[R1RoleParams[arch.UImm8]] {
			return Ldpte(rnd)
		}),
		newInstrCase("Invtlb", func() ohsnap.Arbitrary[U5RRParams] { return Invtlb(rnd) }),
	}
}

// TestInstrGenValid — a generated instruction always encodes into a
// single 4-byte word (LA64: no compression), at the family base pc.
func TestInstrGenValid(t *testing.T) {
	rnd := arb.Rnd(42)
	for _, c := range instrCases(rnd) {
		var buf bytes.Buffer
		for range 300 {
			buf.Reset()
			in, _ := c.sample()
			_, err := in.Encode(&buf, 0)
			require.NoError(t, err, "%s (%s)", c.name, in.ObjDump(disasm.DefaultViewCtx()))
			require.Equal(t, 4, buf.Len(), "%s: %d bytes", c.name, buf.Len())
		}
	}
}

// TestInstrGenText — String() of the parameters is the ObjDump text of
// the instruction (the property suite reports failures by it).
func TestInstrGenText(t *testing.T) {
	rnd := arb.Rnd(42)
	for _, c := range instrCases(rnd) {
		in, s := c.sample()
		require.Equal(
			t,
			in.ObjDump(disasm.DefaultViewCtx()),
			s,
			"%s: params text",
			c.name,
		)
	}
}

// TestInstrGenShrinkValid — shrink candidates are valid instructions
// that still encode (branches stay encodable at BranchBase).
func TestInstrGenShrinkValid(t *testing.T) {
	rnd := arb.Rnd(9)
	for _, c := range instrCases(rnd) {
		var buf bytes.Buffer
		for range 200 {
			buf.Reset()
			for _, in := range c.shrink() {
				_, err := in.Encode(&buf, 0)
				require.NoError(t, err, "%s: shrink candidate %s", c.name,
					in.ObjDump(disasm.DefaultViewCtx()))
			}
		}
	}
}

// --- branch targets --------------------------------------------------------

// TestBranchTargetsProperty — offsets are word-aligned values inside
// each family's span; the shrink skips the halved values that lost
// alignment.
func TestBranchTargetsProperty(t *testing.T) {
	rnd := arb.Rnd(42)
	aligned := func(t int64) bool { return t%4 == 0 }
	ohsnap.Check(t, 300, Branch2(rnd), func(p Branch2Params) bool {
		return aligned(p.Target) &&
			p.Target >= -branch2Span && p.Target <= branch2Span
	})
	ohsnap.Check(t, 300, Branch1(rnd), func(p Branch1Params) bool {
		return aligned(p.Target) &&
			p.Target >= -branch1Span && p.Target <= branch1Span
	})
	ohsnap.Check(t, 300, Jump(rnd), func(p JumpParams) bool {
		return aligned(p.Target) &&
			p.Target >= -jumpSpan && p.Target <= jumpSpan
	})

	// offset 12 halves to 6/3/1/0 — only 0 survives the alignment check
	require.Equal(t, []int64{0}, branchShrunk(12))
	require.Empty(t, branchShrunk(0), "shrink of zero itself")
}

// --- bit-field domain ------------------------------------------------------

// TestFieldParamsProperty — msb >= lsb always holds, and every shrink
// candidate keeps the cross-check (an invalid pair would not assemble).
func TestFieldParamsProperty(t *testing.T) {
	rnd := arb.Rnd(42)
	pair := func(msb, lsb int64) bool { return msb >= lsb }
	ohsnap.Check(t, 300, FieldW(rnd), func(p FieldWParams) bool {
		if !pair(p.Msb.Val(), p.Lsb.Val()) {
			return false
		}

		for c := range FieldW(rnd).Shrink(p) {
			if !pair(c.Msb.Val(), c.Lsb.Val()) {
				return false
			}
		}

		return true
	})
	ohsnap.Check(t, 300, FieldD(rnd), func(p FieldDParams) bool {
		if !pair(p.Msb.Val(), p.Lsb.Val()) {
			return false
		}

		for c := range FieldD(rnd).Shrink(p) {
			if !pair(c.Msb.Val(), c.Lsb.Val()) {
				return false
			}
		}

		return true
	})
}

// --- inventory -------------------------------------------------------------

// TestFamilyCoverage — the family tables cover every mnemonic of the
// arch decode table except the operandless forms (emptyForms: a single
// fixed word each, exercised by the decodeTable-driven differential).
func TestFamilyCoverage(t *testing.T) {
	covered := map[string]bool{}
	for _, e := range alu3R {
		covered[e.name] = true
	}

	for _, e := range alu2R {
		covered[e.name] = true
	}

	for _, e := range aluImm12 {
		covered[e.name] = true
	}

	for _, e := range aluUImm12 {
		covered[e.name] = true
	}

	for _, e := range aluImm16 {
		covered[e.name] = true
	}

	for _, e := range imm20 {
		covered[e.name] = true
	}

	for _, e := range code15 {
		covered[e.name] = true
	}

	for _, e := range branch2 {
		covered[e.name] = true
	}

	for _, e := range branch1 {
		covered[e.name] = true
	}

	for _, e := range jump {
		covered[e.name] = true
	}

	for _, e := range jirl {
		covered[e.name] = true
	}

	for _, e := range ldSt {
		covered[e.name] = true
	}

	for _, e := range ldptr {
		covered[e.name] = true
	}

	for _, e := range ldxStx {
		covered[e.name] = true
	}

	for _, e := range ldAcq {
		covered[e.name] = true
	}

	for _, e := range hints {
		covered[e.name] = true
	}

	for _, e := range preldx {
		covered[e.name] = true
	}

	for _, e := range shiftW {
		covered[e.name] = true
	}

	for _, e := range shiftD {
		covered[e.name] = true
	}

	for _, e := range fieldW {
		covered[e.name] = true
	}

	for _, e := range fieldD {
		covered[e.name] = true
	}

	for _, e := range alsl {
		covered[e.name] = true
	}

	for _, e := range bytepickW {
		covered[e.name] = true
	}

	for _, e := range bytepickD {
		covered[e.name] = true
	}

	for _, e := range atomics {
		covered[e.name] = true
	}

	for _, e := range csrRW {
		covered[e.name] = true
	}

	for _, e := range csrXchg {
		covered[e.name] = true
	}

	for _, e := range ioCsr {
		covered[e.name] = true
	}

	for _, e := range lddir {
		covered[e.name] = true
	}

	for _, e := range ldpte {
		covered[e.name] = true
	}

	for _, e := range invtlb {
		covered[e.name] = true
	}

	empty := map[string]bool{}
	for _, n := range EmptyForms() {
		empty[n] = true
	}

	for _, m := range arch.Mnemonics() {
		if empty[m] {
			require.False(t, covered[m], "%s: an operandless form has no family", m)
			continue
		}

		require.True(t, covered[m], "%s: not covered by any family table", m)
	}

	require.Len(t, empty, len(emptyForms), "EmptyForms returns a copy")
}
