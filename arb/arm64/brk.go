package arm64

// Generator for brk — one generator, one type, one constructor (arm64.NewBrk).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BrkParams — parameters of brk #imm16.
type BrkParams struct {
	Imm arm64.Imm16
}

func NewBrkParams(imm arm64.Imm16) BrkParams {
	return BrkParams{Imm: imm}
}

func (p BrkParams) Instr() arm64.Instr {
	return arm64.NewBrk(p.Imm)
}
func (p BrkParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// brkGen — generator for brk: immediate uniform in 0..0xffff.
type brkGen struct {
	rnd *rand.Rand
}

func newBrkGen(rnd *rand.Rand) brkGen {
	return brkGen{rnd: rnd}
}

// Brk — an arbitrary brk.
func Brk(rnd *rand.Rand) ohsnap.Arbitrary[BrkParams] {
	return newBrkGen(rnd)
}

func (g brkGen) Generate() BrkParams {
	return NewBrkParams(imm16(g.rnd.Int64N(0x10000)))
}

func (g brkGen) Shrink(p BrkParams) []BrkParams {
	v, err := immValue(p.Imm)
	if err != nil {
		return nil // String() of our own type is unparseable — invariant
	}

	var out []BrkParams
	for _, d := range arb.Halved(v) {
		imm, err := arm64.NewImm16(d)
		if err != nil {
			continue // unreachable: half of a valid imm16 is always in 0..65535
		}

		out = append(out, NewBrkParams(imm))
	}

	return out
}
