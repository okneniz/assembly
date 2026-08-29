package arm64

// Generator for add (shifted register) — one generator, one type, one
// constructor (arm64.NewAddShift).

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

// AddShiftParams — parameters of add rd, rn, rm[, shift #imm6].
type AddShiftParams struct {
	Rd, Rn, Rm arm64.Reg
	Imm        arm64.Imm6
	Sh         arm64.Shift
}

func NewAddShiftParams(
	rd arm64.Reg,
	rn arm64.Reg,
	rm arm64.Reg,
	imm arm64.Imm6,
	sh arm64.Shift,
) AddShiftParams {
	return AddShiftParams{
		Rd:  rd,
		Rn:  rn,
		Rm:  rm,
		Imm: imm,
		Sh:  sh,
	}
}

func (p AddShiftParams) Instr() arm64.Instr {
	in, err := arm64.NewAddShift(p.Rd, p.Rn, p.Rm, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddShiftParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// addShiftGen — generator for add: registers of the same width (31st = zr),
// shift lsl/lsr/asr (ror is unallocated for add/sub), amount 0..63 (0..31
// in the 32-bit form).
type addShiftGen struct {
	rnd *rand.Rand
}

func newAddShiftGen(rnd *rand.Rand) addShiftGen {
	return addShiftGen{rnd: rnd}
}

// AddShift — an arbitrary add (shifted register).
func AddShift(rnd *rand.Rand) ohsnap.Arbitrary[AddShiftParams] {
	return newAddShiftGen(rnd)
}

func (g addShiftGen) Generate() iter.Seq[AddShiftParams] {
	return arb.Stream(func() AddShiftParams {
		is64 := g.rnd.IntN(2) == 1
		hi := int64(63)
		if !is64 {
			hi = 31
		}

		return NewAddShiftParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			imm6(g.rnd.Int64N(hi+1)),
			arm64.Shift(g.rnd.IntN(3)),
		)
	})
}

func (g addShiftGen) Shrink(p AddShiftParams) iter.Seq[AddShiftParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[AddShiftParams]() // String() of our own type is unparseable — invariant
	}

	var out []AddShiftParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewAddShiftParams(r, p.Rn, p.Rm, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewAddShiftParams(p.Rd, r, p.Rm, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewAddShiftParams(p.Rd, p.Rn, r, p.Imm, p.Sh))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.NewImm6(d)
		if err != nil {
			continue // unreachable: half of a valid imm6 is always in 0..63
		}

		out = append(out, NewAddShiftParams(p.Rd, p.Rn, p.Rm, imm, p.Sh))
	}

	if p.Sh != arm64.LSL {
		out = append(out, NewAddShiftParams(p.Rd, p.Rn, p.Rm, p.Imm, arm64.LSL))
	}

	return slices.Values(out)
}
