package arm64

// Generator for umulh — one generator, one type, one constructor (Umulh).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UmulhParams — parameters of umulh rd, rn, rm (the 64-bit form only).
type UmulhParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewUmulhParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) UmulhParams {
	return UmulhParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p UmulhParams) Instr() arm64.Instr {
	in, err := arm64.New().Umulh(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p UmulhParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// umulhGen — generator for umulh: x-registers, occasionally xzr.
type umulhGen struct {
	rnd *rand.Rand
}

func newUmulhGen(rnd *rand.Rand) umulhGen {
	return umulhGen{rnd: rnd}
}

// Umulh — an arbitrary umulh.
func Umulh(rnd *rand.Rand) ohsnap.Arbitrary[UmulhParams] {
	return newUmulhGen(rnd)
}

func (g umulhGen) Generate() iter.Seq[UmulhParams] {
	return arb.Stream(func() UmulhParams {
		return NewUmulhParams(
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
		)
	})
}

func (g umulhGen) Shrink(p UmulhParams) iter.Seq[UmulhParams] {
	var out []UmulhParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewUmulhParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewUmulhParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewUmulhParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
