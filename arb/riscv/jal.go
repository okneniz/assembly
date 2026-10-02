package riscv

// Generator for jal - one generator, one type, one constructor (Jal):
// rd and a pc-relative byte offset (even, ±1MB - the generator
// invariant, the Builder stores the offset unvalidated).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// the pc-relative offset bounds of the jal class (even offsets).
const (
	jalOffMin = -1 << 20
	jalOffMax = 1<<20 - 2
)

// JalParams — parameters of jal.
type JalParams struct {
	Rd  riscv.Reg
	Off int64
}

func NewJalParams(rd riscv.Reg, off int64) JalParams {
	return JalParams{
		Rd:  rd,
		Off: off,
	}
}

func (p JalParams) Instr() riscv.Instr {
	return riscv.New().Jal(p.Rd, p.Off)
}

func (p JalParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Jal — an arbitrary jal.
func Jal(rnd *rand.Rand) ohsnap.Arbitrary[JalParams] {
	return jalGen{rnd: rnd}
}

type jalGen struct {
	rnd *rand.Rand
}

func (g jalGen) Generate() iter.Seq[JalParams] {
	return arb.Stream(func() JalParams {
		return NewJalParams(reg(g.rnd), jalOffMin+2*g.rnd.Int64N((jalOffMax-jalOffMin)/2+1))
	})
}

func (g jalGen) Shrink(p JalParams) iter.Seq[JalParams] {
	out := make([]JalParams, 0, 8)
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewJalParams(r, p.Off))
	}

	for _, v := range pcOffShrunk(p.Off, jalOffMin, jalOffMax) {
		out = append(out, NewJalParams(p.Rd, v))
	}

	return slices.Values(out)
}
