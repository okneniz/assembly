package arm64

// Generator for asr (register) — one generator, one type, one constructor
// (AsrReg).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AsrRegParams — parameters of asr rd, rn, rm (register).
type AsrRegParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewAsrRegParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) AsrRegParams {
	return AsrRegParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p AsrRegParams) Instr() arm64.Instr {
	in, err := arm64.New().AsrReg(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p AsrRegParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// asrRegGen — generator for asr: registers of the same width, 31st is zr.
type asrRegGen struct {
	rnd *rand.Rand
}

// AsrReg — an arbitrary asr (register).
func AsrReg(rnd *rand.Rand) ohsnap.Arbitrary[AsrRegParams] {
	return newAsrRegGen(rnd)
}

func newAsrRegGen(rnd *rand.Rand) asrRegGen {
	return asrRegGen{rnd: rnd}
}

func (g asrRegGen) Generate() iter.Seq[AsrRegParams] {
	return arb.Stream(func() AsrRegParams {
		is64 := g.rnd.IntN(2) == 1
		return NewAsrRegParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
		)
	})
}

func (g asrRegGen) Shrink(p AsrRegParams) iter.Seq[AsrRegParams] {
	regs := regShrunk(p.Rd)
	out := make([]AsrRegParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewAsrRegParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewAsrRegParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewAsrRegParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
