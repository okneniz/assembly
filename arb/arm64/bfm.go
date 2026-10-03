package arm64

// Generator for bfm — one generator, one type, one constructor (Bfm).
// The plain BFM is general (any immr/imms below the register width); the
// ctor does not check the width bound, the generator does.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BfmParams — parameters of bfm rd, rn, #immr, #imms.
type BfmParams struct {
	Rd, Rn     arm64.Reg // 31 reads as zr
	Immr, Imms uint32
}

func NewBfmParams(rd arm64.Reg, rn arm64.Reg, immr uint32, imms uint32) BfmParams {
	return BfmParams{
		Rd:   rd,
		Rn:   rn,
		Immr: immr,
		Imms: imms,
	}
}

func (p BfmParams) Instr() arm64.Instr {
	in, err := arm64.New().Bfm(p.Rd, p.Rn, p.Immr, p.Imms)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p BfmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// bfmGen — generator for bfm: registers of one width, immr/imms below
// the register width.
type bfmGen struct {
	rnd *rand.Rand
}

// Bfm — an arbitrary bfm.
func Bfm(rnd *rand.Rand) ohsnap.Arbitrary[BfmParams] {
	return newBfmGen(rnd)
}

func newBfmGen(rnd *rand.Rand) bfmGen {
	return bfmGen{rnd: rnd}
}

func (g bfmGen) Generate() iter.Seq[BfmParams] {
	return arbStream(func() BfmParams {
		hi := uint32(64)
		is64 := g.rnd.IntN(2) == 1
		if !is64 {
			hi = 32
		}

		return NewBfmParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			uint32(g.rnd.IntN(int(hi))),
			uint32(g.rnd.IntN(int(hi))),
		)
	})
}

func (g bfmGen) Shrink(p BfmParams) iter.Seq[BfmParams] {
	regs := regShrunk(p.Rd)
	out := make([]BfmParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewBfmParams(r, p.Rn, p.Immr, p.Imms))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewBfmParams(p.Rd, r, p.Immr, p.Imms))
	}

	for _, v := range u32Halved(p.Immr) {
		out = append(out, NewBfmParams(p.Rd, p.Rn, v, p.Imms))
	}

	for _, v := range u32Halved(p.Imms) {
		out = append(out, NewBfmParams(p.Rd, p.Rn, p.Immr, v))
	}

	return slices.Values(out)
}
