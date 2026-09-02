package arm64

// Generator for sub (shifted register) — one generator, one type, one
// constructor (SubShift).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SubShiftParams — parameters of sub rd, rn, rm[, shift #imm6].
type SubShiftParams struct {
	Rd, Rn, Rm arm64.Reg
	Imm        arm64.Imm6
	Sh         arm64.Shift
}

func NewSubShiftParams(
	rd arm64.Reg,
	rn arm64.Reg,
	rm arm64.Reg,
	imm arm64.Imm6,
	sh arm64.Shift,
) SubShiftParams {
	return SubShiftParams{
		Rd:  rd,
		Rn:  rn,
		Rm:  rm,
		Imm: imm,
		Sh:  sh,
	}
}

func (p SubShiftParams) Instr() arm64.Instr {
	in, err := arm64.New().SubShift(p.Rd, p.Rn, p.Rm, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SubShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// subShiftGen — generator for sub: registers of the same width (31st = zr),
// shift lsl/lsr/asr (ror is unallocated for add/sub), amount 0..63 (0..31
// in the 32-bit form).
type subShiftGen struct {
	rnd *rand.Rand
}

func newSubShiftGen(rnd *rand.Rand) subShiftGen {
	return subShiftGen{rnd: rnd}
}

// SubShift — an arbitrary sub (shifted register).
func SubShift(rnd *rand.Rand) ohsnap.Arbitrary[SubShiftParams] {
	return newSubShiftGen(rnd)
}

func (g subShiftGen) Generate() iter.Seq[SubShiftParams] {
	return arb.Stream(func() SubShiftParams {
		is64 := g.rnd.IntN(2) == 1
		hi := int64(63)
		if !is64 {
			hi = 31
		}

		return NewSubShiftParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			imm6(g.rnd.Int64N(hi+1)),
			arm64.Shift(g.rnd.IntN(3)),
		)
	})
}

func (g subShiftGen) Shrink(p SubShiftParams) iter.Seq[SubShiftParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[SubShiftParams]() // String() of our own type is unparseable — invariant
	}

	var out []SubShiftParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewSubShiftParams(r, p.Rn, p.Rm, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewSubShiftParams(p.Rd, r, p.Rm, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewSubShiftParams(p.Rd, p.Rn, r, p.Imm, p.Sh))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm6(d)
		if err != nil {
			continue // unreachable: half of a valid imm6 is always in 0..63
		}

		out = append(out, NewSubShiftParams(p.Rd, p.Rn, p.Rm, imm, p.Sh))
	}

	if p.Sh != arm64.LSL {
		out = append(out, NewSubShiftParams(p.Rd, p.Rn, p.Rm, p.Imm, arm64.LSL))
	}

	return slices.Values(out)
}
