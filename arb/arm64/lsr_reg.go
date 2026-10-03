package arm64

// Generator for lsr (register) — one generator, one type, one constructor
// (LsrReg).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LsrRegParams — parameters of lsr rd, rn, rm (register).
type LsrRegParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewLsrRegParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) LsrRegParams {
	return LsrRegParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p LsrRegParams) Instr() arm64.Instr {
	in, err := arm64.New().LsrReg(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p LsrRegParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// lsrRegGen — generator for lsr: registers of the same width, 31st is zr.
type lsrRegGen struct {
	rnd *rand.Rand
}

// LsrReg — an arbitrary lsr (register).
func LsrReg(rnd *rand.Rand) ohsnap.Arbitrary[LsrRegParams] {
	return newLsrRegGen(rnd)
}

func newLsrRegGen(rnd *rand.Rand) lsrRegGen {
	return lsrRegGen{rnd: rnd}
}

func (g lsrRegGen) Generate() iter.Seq[LsrRegParams] {
	return arb.Stream(func() LsrRegParams {
		is64 := g.rnd.IntN(2) == 1
		return NewLsrRegParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
		)
	})
}

func (g lsrRegGen) Shrink(p LsrRegParams) iter.Seq[LsrRegParams] {
	regs := regShrunk(p.Rd)
	out := make([]LsrRegParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewLsrRegParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewLsrRegParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewLsrRegParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
