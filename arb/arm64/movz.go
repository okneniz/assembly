package arm64

// Generator for movz — one generator, one type, one constructor (Movz).

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

// MovzParams — parameters of movz rd, #imm16, lsl #hw*16.
type MovzParams struct {
	Rd  arm64.Reg
	Imm arm64.Imm16
	Hw  arm64.Hw
}

func NewMovzParams(rd arm64.Reg, imm arm64.Imm16, hw arm64.Hw) MovzParams {
	return MovzParams{
		Rd:  rd,
		Imm: imm,
		Hw:  hw,
	}
}

func (p MovzParams) Instr() arm64.Instr {
	in, err := arm64.New().Movz(p.Rd, p.Imm, p.Hw)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p MovzParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// movzGen — generator for movz: the class of rd is consistent with the hw
// shift (the 32-bit form allows only Hw0/Hw1), sp is not allowed.
type movzGen struct {
	rnd *rand.Rand
}

func newMovzGen(rnd *rand.Rand) movzGen {
	return movzGen{rnd: rnd}
}

// Movz — an arbitrary movz.
func Movz(rnd *rand.Rand) ohsnap.Arbitrary[MovzParams] {
	return newMovzGen(rnd)
}

func (g movzGen) Generate() iter.Seq[MovzParams] {
	return arb.Stream(func() MovzParams {
		hw := arm64.Hw(g.rnd.IntN(4))
		// hw>=Hw2 (lsl #32/#48) occurs only in the 64-bit form.
		is64 := hw >= arm64.Hw2 || g.rnd.IntN(2) == 1
		return NewMovzParams(genReg(g.rnd, is64, false, true), imm16(g.rnd.Int64N(0x10000)), hw)
	})
}

func (g movzGen) Shrink(p MovzParams) iter.Seq[MovzParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnap.Empty[MovzParams]() // String() of our own type is unparseable — invariant
	}

	var out []MovzParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewMovzParams(r, p.Imm, p.Hw))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm16(d)
		if err != nil {
			continue // unreachable: half of a valid imm16 is always in 0..65535
		}

		out = append(out, NewMovzParams(p.Rd, imm, p.Hw))
	}

	if p.Hw != arm64.Hw0 {
		out = append(out, NewMovzParams(p.Rd, p.Imm, arm64.Hw0))
	}

	return slices.Values(out)
}
