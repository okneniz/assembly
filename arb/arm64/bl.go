package arm64

// Generator for bl — one generator, one type, one constructor (Bl).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BlParams — parameters of bl (the pc-relative byte offset).
type BlParams struct {
	Off int64
}

func NewBlParams(off int64) BlParams {
	return BlParams{Off: off}
}

func (p BlParams) Instr() arm64.Instr {
	return arm64.New().Bl(p.Off)
}
func (p BlParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// blGen — generator for bl: the offset uniform in the ±128MB imm26 range.
type blGen struct {
	off ohsnap.Arbitrary[int64]
}

func newBlGen(rnd *rand.Rand) blGen {
	return blGen{off: BrOff(rnd, 1<<27)}
}

// Bl — an arbitrary bl.
func Bl(rnd *rand.Rand) ohsnap.Arbitrary[BlParams] {
	return newBlGen(rnd)
}

func (g blGen) Generate() iter.Seq[BlParams] {
	return arb.Stream(func() BlParams {
		return NewBlParams(ohsnap.First(g.off.Generate()))
	})
}

func (g blGen) Shrink(p BlParams) iter.Seq[BlParams] {
	var out []BlParams
	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewBlParams(v))
	}

	return slices.Values(out)
}
