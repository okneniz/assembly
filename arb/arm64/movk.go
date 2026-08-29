package arm64

// Generator for movk — one generator, one type, one constructor (arm64.NewMovk).

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

// MovkParams — parameters of movk rd, #imm16, lsl #hw*16.
type MovkParams struct {
	Rd  arm64.Reg
	Imm arm64.Imm16
	Hw  arm64.Hw
}

func NewMovkParams(rd arm64.Reg, imm arm64.Imm16, hw arm64.Hw) MovkParams {
	return MovkParams{
		Rd:  rd,
		Imm: imm,
		Hw:  hw,
	}
}

func (p MovkParams) Instr() arm64.Instr {
	in, err := arm64.NewMovk(p.Rd, p.Imm, p.Hw)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p MovkParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// movkGen — generator for movk: the class of rd is consistent with the hw
// shift (the 32-bit form allows only Hw0/Hw1), sp is not allowed.
type movkGen struct {
	rnd *rand.Rand
}

func newMovkGen(rnd *rand.Rand) movkGen {
	return movkGen{rnd: rnd}
}

// Movk — an arbitrary movk.
func Movk(rnd *rand.Rand) ohsnap.Arbitrary[MovkParams] {
	return newMovkGen(rnd)
}

func (g movkGen) Generate() iter.Seq[MovkParams] {
	return arb.Stream(func() MovkParams {
		hw := arm64.Hw(g.rnd.IntN(4))
		// hw>=Hw2 (lsl #32/#48) occurs only in the 64-bit form.
		is64 := hw >= arm64.Hw2 || g.rnd.IntN(2) == 1
		return NewMovkParams(genReg(g.rnd, is64, false, true), imm16(g.rnd.Int64N(0x10000)), hw)
	})
}

func (g movkGen) Shrink(p MovkParams) iter.Seq[MovkParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[MovkParams]() // String() of our own type is unparseable — invariant
	}

	var out []MovkParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewMovkParams(r, p.Imm, p.Hw))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.NewImm16(d)
		if err != nil {
			continue // unreachable: half of a valid imm16 is always in 0..65535
		}

		out = append(out, NewMovkParams(p.Rd, imm, p.Hw))
	}

	if p.Hw != arm64.Hw0 {
		out = append(out, NewMovkParams(p.Rd, p.Imm, arm64.Hw0))
	}

	return slices.Values(out)
}
