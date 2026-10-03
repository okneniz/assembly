package arm64

// Generator for movn — one generator, one type, one constructor (Movn).

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

// MovnParams — parameters of movn rd, #imm16, lsl #hw*16.
type MovnParams struct {
	Rd  arm64.Reg
	Imm arm64.Imm16
	Hw  arm64.Hw
}

func NewMovnParams(rd arm64.Reg, imm arm64.Imm16, hw arm64.Hw) MovnParams {
	return MovnParams{
		Rd:  rd,
		Imm: imm,
		Hw:  hw,
	}
}

func (p MovnParams) Instr() arm64.Instr {
	in, err := arm64.New().Movn(p.Rd, p.Imm, p.Hw)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p MovnParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// movnGen — generator for movn: the class of rd is consistent with the hw
// shift (the 32-bit form allows only Hw0/Hw1), sp is not allowed.
type movnGen struct {
	rnd *rand.Rand
}

// Movn — an arbitrary movn.
func Movn(rnd *rand.Rand) ohsnap.Arbitrary[MovnParams] {
	return newMovnGen(rnd)
}

func newMovnGen(rnd *rand.Rand) movnGen {
	return movnGen{rnd: rnd}
}

func (g movnGen) Generate() iter.Seq[MovnParams] {
	return arb.Stream(func() MovnParams {
		hw := arm64.Hw(g.rnd.IntN(4))
		// hw>=Hw2 (lsl #32/#48) occurs only in the 64-bit form.
		is64 := hw >= arm64.Hw2 || g.rnd.IntN(2) == 1
		return NewMovnParams(genReg(g.rnd, is64, false, true), imm16(g.rnd.Int64N(0x10000)), hw)
	})
}

func (g movnGen) Shrink(p MovnParams) iter.Seq[MovnParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[MovnParams]() // String() of our own type is unparseable — invariant
	}

	var out []MovnParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewMovnParams(r, p.Imm, p.Hw))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm16(d)
		if err != nil {
			continue // unreachable: half of a valid imm16 is always in 0..65535
		}

		out = append(out, NewMovnParams(p.Rd, imm, p.Hw))
	}

	if p.Hw != arm64.Hw0 {
		out = append(out, NewMovnParams(p.Rd, p.Imm, arm64.Hw0))
	}

	return slices.Values(out)
}
