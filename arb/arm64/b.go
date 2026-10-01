package arm64

// Generator for b — one generator, one type, one constructor (B).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BParams — parameters of b (the pc-relative byte offset).
type BParams struct {
	Off int64
}

func NewBParams(off int64) BParams {
	return BParams{Off: off}
}

func (p BParams) Instr() arm64.Instr {
	return arm64.New().B(p.Off)
}
func (p BParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// bGen — generator for b: the offset uniform in the ±128MB imm26 range.
type bGen struct {
	off ohsnap.Arbitrary[int64]
}

func newBGen(rnd *rand.Rand) bGen {
	return bGen{off: BrOff(rnd, 1<<27)}
}

// B — an arbitrary b.
func B(rnd *rand.Rand) ohsnap.Arbitrary[BParams] {
	return newBGen(rnd)
}

func (g bGen) Generate() iter.Seq[BParams] {
	return arb.Stream(func() BParams {
		return NewBParams(ohsnap.First(g.off.Generate()))
	})
}

func (g bGen) Shrink(p BParams) iter.Seq[BParams] {
	var out []BParams
	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewBParams(v))
	}

	return slices.Values(out)
}
