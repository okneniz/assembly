package arm64

// The scalar FP family cores: the same-kind forms (one s or d file
// through all operands — requireFpKind), the width-crossing conversions
// (fcvt: opposite kinds; gpr↔fpr: any pairing), and the immediate form
// (the VfpExpand imm8 space). The per-instruction files are thin
// constructors over these.

import (
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/assembly/arch/arm64"
)

// F2Params — rd, rn of one FP kind.
type F2Params struct {
	Rd, Rn arm64.FReg
}

func NewF2Params(rd arm64.FReg, rn arm64.FReg) F2Params {
	return F2Params{
		Rd: rd,
		Rn: rn,
	}
}

// f2Gen — the two-register same-kind core.
type f2Gen struct {
	rnd *rand.Rand
}

func (g f2Gen) Generate() iter.Seq[F2Params] {
	return arbStream(func() F2Params {
		is64 := g.rnd.IntN(2) == 1
		return NewF2Params(genFp(g.rnd, is64), genFp(g.rnd, is64))
	})
}

func (g f2Gen) Shrink(p F2Params) iter.Seq[F2Params] {
	fps := fpShrunk(p.Rd)
	out := make([]F2Params, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewF2Params(r, p.Rn))
	}

	for _, r := range fpShrunk(p.Rn) {
		out = append(out, NewF2Params(p.Rd, r))
	}

	return slices.Values(out)
}

// F3Params — rd, rn, rm of one FP kind.
type F3Params struct {
	Rd, Rn, Rm arm64.FReg
}

func NewF3Params(rd arm64.FReg, rn arm64.FReg, rm arm64.FReg) F3Params {
	return F3Params{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

// f3Gen — the three-register same-kind core.
type f3Gen struct {
	rnd *rand.Rand
}

func (g f3Gen) Generate() iter.Seq[F3Params] {
	return arbStream(func() F3Params {
		is64 := g.rnd.IntN(2) == 1
		return NewF3Params(genFp(g.rnd, is64), genFp(g.rnd, is64), genFp(g.rnd, is64))
	})
}

func (g f3Gen) Shrink(p F3Params) iter.Seq[F3Params] {
	fps := fpShrunk(p.Rd)
	out := make([]F3Params, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewF3Params(r, p.Rn, p.Rm))
	}

	for _, r := range fpShrunk(p.Rn) {
		out = append(out, NewF3Params(p.Rd, r, p.Rm))
	}

	for _, r := range fpShrunk(p.Rm) {
		out = append(out, NewF3Params(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}

// F4Params — rd, rn, rm, ra of one FP kind.
type F4Params struct {
	Rd, Rn, Rm, Ra arm64.FReg
}

func NewF4Params(rd arm64.FReg, rn arm64.FReg, rm arm64.FReg, ra arm64.FReg) F4Params {
	return F4Params{
		Rd: rd,
		Rn: rn,
		Rm: rm,
		Ra: ra,
	}
}

// f4Gen — the four-register same-kind core.
type f4Gen struct {
	rnd *rand.Rand
}

func (g f4Gen) Generate() iter.Seq[F4Params] {
	return arbStream(func() F4Params {
		is64 := g.rnd.IntN(2) == 1
		return NewF4Params(
			genFp(g.rnd, is64),
			genFp(g.rnd, is64),
			genFp(g.rnd, is64),
			genFp(g.rnd, is64),
		)
	})
}

func (g f4Gen) Shrink(p F4Params) iter.Seq[F4Params] {
	fps := fpShrunk(p.Rd)
	out := make([]F4Params, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewF4Params(r, p.Rn, p.Rm, p.Ra))
	}

	for _, r := range fpShrunk(p.Rn) {
		out = append(out, NewF4Params(p.Rd, r, p.Rm, p.Ra))
	}

	for _, r := range fpShrunk(p.Rm) {
		out = append(out, NewF4Params(p.Rd, p.Rn, r, p.Ra))
	}

	for _, r := range fpShrunk(p.Ra) {
		out = append(out, NewF4Params(p.Rd, p.Rn, p.Rm, r))
	}

	return slices.Values(out)
}

// FGprParams — an FP register and a gpr of the conversion forms (the
// pairing w↔s / x↔d for fmov, any pairing for the cvt family).
type FGprParams struct {
	F   arm64.FReg
	R   arm64.Reg // 31 reads as zr
	Any bool      // any gpr↔fpr pairing (the cvt forms)
}

func NewFGprParams(f arm64.FReg, r arm64.Reg, anyPair bool) FGprParams {
	return FGprParams{
		F:   f,
		R:   r,
		Any: anyPair,
	}
}

// fgprGen — the conversion core: the gpr width follows the FP kind
// (fmov), or any pairing (cvt).
type fgprGen struct {
	rnd *rand.Rand
	any bool
}

func (g fgprGen) Generate() iter.Seq[FGprParams] {
	return arbStream(func() FGprParams {
		is64 := g.rnd.IntN(2) == 1
		r64 := is64
		if g.any {
			r64 = g.rnd.IntN(2) == 1
		}

		return NewFGprParams(
			genFp(g.rnd, is64),
			genReg(g.rnd, r64, false, true),
			g.any,
		)
	})
}

func (g fgprGen) Shrink(p FGprParams) iter.Seq[FGprParams] {
	fps := fpShrunk(p.F)
	out := make([]FGprParams, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewFGprParams(r, p.R, p.Any))
	}

	for _, r := range regShrunk(p.R) {
		out = append(out, NewFGprParams(p.F, r, p.Any))
	}

	return slices.Values(out)
}

// fpImmGen — the immediate core: a random imm8 expanded by the arch
// VfpExpand tables (every imm8 is encodable by construction).
type fpImmGen struct {
	rnd *rand.Rand
}

// FpImmParams — the destination and the imm8 of the value.
type FpImmParams struct {
	Rd  arm64.FReg
	Imm uint8
}

func NewFpImmParams(rd arm64.FReg, imm uint8) FpImmParams {
	return FpImmParams{
		Rd:  rd,
		Imm: imm,
	}
}

// genFp — an FP register of the given kind (s0..s31 / d0..d31).
func genFp(rnd *rand.Rand, is64 bool) arm64.FReg {
	n := rnd.IntN(32)
	if is64 {
		return mustD(n)
	}

	return mustS(n)
}

// fpShrunk — the shrink candidates of an FP register (s0/d0 of the same
// kind, then the halved number).
func fpShrunk(r arm64.FReg) []arm64.FReg {
	var out []arm64.FReg
	zero := mustS(0)
	if r.Is64() {
		zero = mustD(0)
	}

	if r != zero {
		out = append(out, zero)
	}

	n := int(r.Num())
	if n > 1 {
		if r.Is64() {
			out = append(out, mustD(n/2))
		} else {
			out = append(out, mustS(n/2))
		}
	}

	return out
}

func newF2Gen(rnd *rand.Rand) f2Gen {
	return f2Gen{rnd: rnd}
}

func f2(rnd *rand.Rand) f2Gen {
	return newF2Gen(rnd)
}

func newF3Gen(rnd *rand.Rand) f3Gen {
	return f3Gen{rnd: rnd}
}

func f3(rnd *rand.Rand) f3Gen {
	return newF3Gen(rnd)
}

func newF4Gen(rnd *rand.Rand) f4Gen {
	return f4Gen{rnd: rnd}
}

func f4(rnd *rand.Rand) f4Gen {
	return newF4Gen(rnd)
}

func newFGprGen(rnd *rand.Rand, anyPair bool) fgprGen {
	return fgprGen{
		rnd: rnd,
		any: anyPair,
	}
}

// fgpr — the paired core (fmov gpr forms).
func fgpr(rnd *rand.Rand) fgprGen {
	return newFGprGen(rnd, false)
}

// fgprAny — the any-pairing core (the cvt forms).
func fgprAny(rnd *rand.Rand) fgprGen {
	return newFGprGen(rnd, true)
}

func newFpImmGen(rnd *rand.Rand) fpImmGen {
	return fpImmGen{rnd: rnd}
}

// Value — the expandable float of the destination's kind.
func (p FpImmParams) Value() float64 {
	if p.Rd.Is64() {
		return arm64.VfpExpandImm64(uint32(p.Imm))
	}

	return float64(arm64.VfpExpandImm32(uint32(p.Imm)))
}

func (g fpImmGen) Generate() iter.Seq[FpImmParams] {
	return arbStream(func() FpImmParams {
		return NewFpImmParams(genFp(g.rnd, g.rnd.IntN(2) == 1), uint8(g.rnd.IntN(256)))
	})
}

func (g fpImmGen) Shrink(p FpImmParams) iter.Seq[FpImmParams] {
	fps := fpShrunk(p.Rd)
	out := make([]FpImmParams, 0, len(fps))
	for _, r := range fps {
		out = append(out, NewFpImmParams(r, p.Imm))
	}

	for _, v := range u32Halved(uint32(p.Imm)) {
		out = append(out, NewFpImmParams(p.Rd, uint8(v)))
	}

	return slices.Values(out)
}
