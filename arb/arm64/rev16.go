package arm64

// Generator for rev16 — one generator, one type, one constructor (Rev16).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// Rev16Params — parameters of rev16 rd, rn.
type Rev16Params struct {
	Rd, Rn arm64.Reg
}

func NewRev16Params(rd arm64.Reg, rn arm64.Reg) Rev16Params {
	return Rev16Params{
		Rd: rd,
		Rn: rn,
	}
}

func (p Rev16Params) Instr() arm64.Instr {
	in, err := arm64.New().Rev16(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p Rev16Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// rev16Gen — generator for rev16: registers of the same width, 31st is zr.
type rev16Gen struct {
	rnd *rand.Rand
}

// Rev16 — an arbitrary rev16.
func Rev16(rnd *rand.Rand) ohsnap.Arbitrary[Rev16Params] {
	return newRev16Gen(rnd)
}

func newRev16Gen(rnd *rand.Rand) rev16Gen {
	return rev16Gen{rnd: rnd}
}

func (g rev16Gen) Generate() iter.Seq[Rev16Params] {
	return arb.Stream(func() Rev16Params {
		is64 := g.rnd.IntN(2) == 1
		return NewRev16Params(genReg(g.rnd, is64, false, true), genReg(g.rnd, is64, false, true))
	})
}

func (g rev16Gen) Shrink(p Rev16Params) iter.Seq[Rev16Params] {
	regs := regShrunk(p.Rd)
	out := make([]Rev16Params, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewRev16Params(r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewRev16Params(p.Rd, r))
	}

	return slices.Values(out)
}
