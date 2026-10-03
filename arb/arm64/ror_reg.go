package arm64

// Generator for ror (register) — one generator, one type, one constructor
// (RorReg).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// RorRegParams — parameters of ror rd, rn, rm (register).
type RorRegParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewRorRegParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) RorRegParams {
	return RorRegParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p RorRegParams) Instr() arm64.Instr {
	in, err := arm64.New().RorReg(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p RorRegParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// rorRegGen — generator for ror: registers of the same width, 31st is zr.
type rorRegGen struct {
	rnd *rand.Rand
}

// RorReg — an arbitrary ror (register).
func RorReg(rnd *rand.Rand) ohsnap.Arbitrary[RorRegParams] {
	return newRorRegGen(rnd)
}

func newRorRegGen(rnd *rand.Rand) rorRegGen {
	return rorRegGen{rnd: rnd}
}

func (g rorRegGen) Generate() iter.Seq[RorRegParams] {
	return arb.Stream(func() RorRegParams {
		is64 := g.rnd.IntN(2) == 1
		return NewRorRegParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
		)
	})
}

func (g rorRegGen) Shrink(p RorRegParams) iter.Seq[RorRegParams] {
	regs := regShrunk(p.Rd)
	out := make([]RorRegParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewRorRegParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewRorRegParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewRorRegParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
