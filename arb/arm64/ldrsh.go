package arm64

// Generator for ldrsh — one constructor (Ldrsh) over the shared
// load/store core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrshParams — parameters of the ldrsh form.
type LdrshParams struct {
	LsParams
}

func NewLdrshParams(p LsParams) LdrshParams {
	return LdrshParams{LsParams: p}
}

func (p LdrshParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldrsh(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrshParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Ldrsh — an arbitrary ldrsh.
func Ldrsh(rnd *rand.Rand) ohsnap.Arbitrary[LdrshParams] {
	rtX := true
	base := newLsGen(rnd, &rtX, ScaledOff(rnd, 1))
	return ldrshArb{base: base}
}

type ldrshArb struct {
	base lsGen
}

func (a ldrshArb) Generate() iter.Seq[LdrshParams] {
	return arbStream(func() LdrshParams {
		return NewLdrshParams(ohsnap.First(a.base.Generate()))
	})
}

func (a ldrshArb) Shrink(p LdrshParams) iter.Seq[LdrshParams] {
	var out []LdrshParams
	for _, s := range slices.Collect(a.base.Shrink(p.LsParams)) {
		out = append(out, NewLdrshParams(s))
	}

	return slices.Values(out)
}
