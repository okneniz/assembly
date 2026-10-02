package arm64

// Generator for mov (register, SIMD) — one constructor (MovSimd) over
// the shared two-register vector core (the second operand is rm).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// MovSimdParams — parameters of the SIMD mov (register).
type MovSimdParams struct {
	Rd, Rm arm64.VReg
	Arr    string
}

func NewMovSimdParams(rd arm64.VReg, rm arm64.VReg, arr string) MovSimdParams {
	return MovSimdParams{
		Rd:  rd,
		Rm:  rm,
		Arr: arr,
	}
}

func (p MovSimdParams) Instr() arm64.Instr {
	in, err := arm64.New().MovSimd(p.Rd, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p MovSimdParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// movSimdGen — generator for the SIMD mov (register): the logical set.
type movSimdGen struct {
	base v2Gen
}

// MovSimd — an arbitrary SIMD mov (register).
func MovSimd(rnd *rand.Rand) ohsnap.Arbitrary[MovSimdParams] {
	return movSimdGen{base: newV2Gen(rnd, arrLogical())}
}

func (g movSimdGen) Generate() iter.Seq[MovSimdParams] {
	return arbStream(func() MovSimdParams {
		p := ohsnap.First(g.base.Generate())
		return NewMovSimdParams(p.Rd, p.Rn, p.Arr)
	})
}

func (g movSimdGen) Shrink(p MovSimdParams) iter.Seq[MovSimdParams] {
	var out []MovSimdParams
	for _, s := range slices.Collect(g.base.Shrink(V2Params{
		Rd:  p.Rd,
		Rn:  p.Rm,
		Arr: p.Arr,
	})) {
		out = append(out, NewMovSimdParams(s.Rd, s.Rn, s.Arr))
	}

	return slices.Values(out)
}
