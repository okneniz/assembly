package arm64

// The shifted-register logical family core (and/ands/orr/eor/orn/eon/
// bic/bics rd, rn, rm, shift #imm6): one parameter shape, one generator;
// the eight mnemonics are thin constructors over it. Unlike add/sub, ror
// is a legal shift here; the shift amount stays below the register width.

import (
	"iter"

	"github.com/okneniz/assembly/arch/arm64"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/oh-snap/shrink"
)

// ShiftedParams — parameters of the shifted-register logical forms.
type ShiftedParams struct {
	Rd, Rn, Rm arm64.Reg // 31 reads as zr
	Imm6       arm64.Imm6
	Sh         arm64.Shift
}

func NewShiftedParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg, imm6 arm64.Imm6, sh arm64.Shift) ShiftedParams {
	return ShiftedParams{
		Rd:   rd,
		Rn:   rn,
		Rm:   rm,
		Imm6: imm6,
		Sh:   sh,
	}
}

// shiftedGen — the family generator.
type shiftedGen struct {
	rnd *rand.Rand
}

func newShiftedGen(rnd *rand.Rand) shiftedGen {
	return shiftedGen{rnd: rnd}
}

// shifted — the shared Generate/Shrink core of the eight logical families.
func shifted(rnd *rand.Rand) shiftedGen {
	return newShiftedGen(rnd)
}

func (g shiftedGen) Generate() iter.Seq[ShiftedParams] {
	return arbStream(func() ShiftedParams {
		is64 := g.rnd.IntN(2) == 1
		hi := int64(63)
		if !is64 {
			hi = 31
		}

		return NewShiftedParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			imm6(g.rnd.Int64N(hi+1)),
			arm64.Shift(g.rnd.IntN(4)),
		)
	})
}

func (g shiftedGen) Shrink(p ShiftedParams) iter.Seq[ShiftedParams] {
	v, err := immValue(p.Imm6)
	if err != nil {
		return ohsnapEmpty[ShiftedParams]() // String() of our own type is unparseable — invariant
	}

	var out []ShiftedParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewShiftedParams(r, p.Rn, p.Rm, p.Imm6, p.Sh))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewShiftedParams(p.Rd, r, p.Rm, p.Imm6, p.Sh))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewShiftedParams(p.Rd, p.Rn, r, p.Imm6, p.Sh))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm6(d)
		if err != nil {
			continue // unreachable: half of a valid imm6 is always in 0..63
		}

		out = append(out, NewShiftedParams(p.Rd, p.Rn, p.Rm, imm, p.Sh))
	}

	return slices.Values(out)
}
