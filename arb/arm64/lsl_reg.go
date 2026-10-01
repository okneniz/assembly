package arm64

// Generator for lsl (register) — one generator, one type, one constructor
// (LslReg).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LslRegParams — parameters of lsl rd, rn, rm (register).
type LslRegParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewLslRegParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) LslRegParams {
	return LslRegParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p LslRegParams) Instr() arm64.Instr {
	in, err := arm64.New().LslReg(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LslRegParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// lslRegGen — generator for lsl: registers of the same width, 31st is zr.
type lslRegGen struct {
	rnd *rand.Rand
}

func newLslRegGen(rnd *rand.Rand) lslRegGen {
	return lslRegGen{rnd: rnd}
}

// LslReg — an arbitrary lsl (register).
func LslReg(rnd *rand.Rand) ohsnap.Arbitrary[LslRegParams] {
	return newLslRegGen(rnd)
}

func (g lslRegGen) Generate() iter.Seq[LslRegParams] {
	return arb.Stream(func() LslRegParams {
		is64 := g.rnd.IntN(2) == 1
		return NewLslRegParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
		)
	})
}

func (g lslRegGen) Shrink(p LslRegParams) iter.Seq[LslRegParams] {
	var out []LslRegParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewLslRegParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewLslRegParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewLslRegParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
