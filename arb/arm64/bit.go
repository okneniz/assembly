package arm64

// Generator for bit — one constructor (Bit) over the shared
// three-register vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// BitParams — parameters of bit.
type BitParams struct {
	V3Params
}

func NewBitParams(p V3Params) BitParams {
	return BitParams{V3Params: p}
}

func (p BitParams) Instr() arm64.Instr {
	in, err := arm64.New().Bit(p.Rd, p.Rn, p.Rm, p.Arr)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p BitParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Bit — an arbitrary bit.
func Bit(rnd *rand.Rand) ohsnap.Arbitrary[BitParams] {
	base := newV3Gen(rnd, arrLogical())
	return bitArb{base: base}
}

type bitArb struct {
	base v3Gen
}

func (a bitArb) Generate() iter.Seq[BitParams] {
	return arbStream(func() BitParams {
		return NewBitParams(ohsnap.First(a.base.Generate()))
	})
}

func (a bitArb) Shrink(p BitParams) iter.Seq[BitParams] {
	var out []BitParams
	for _, s := range slices.Collect(a.base.Shrink(p.V3Params)) {
		out = append(out, NewBitParams(s))
	}

	return slices.Values(out)
}
