package arm64

// Generator for tbl — one constructor (Tbl) over the shared
// arrangement-less vector core.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// TblParams — parameters of tbl.
type TblParams struct {
	V3PlainParams
}

func NewTblParams(p V3PlainParams) TblParams {
	return TblParams{V3PlainParams: p}
}

func (p TblParams) Instr() arm64.Instr {
	in, err := arm64.New().Tbl(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p TblParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Tbl — an arbitrary tbl.
func Tbl(rnd *rand.Rand) ohsnap.Arbitrary[TblParams] {
	base := v3Plain(rnd)
	return tblArb{base: base}
}

type tblArb struct {
	base v3PlainGen
}

func (a tblArb) Generate() iter.Seq[TblParams] {
	return arbStream(func() TblParams {
		return NewTblParams(ohsnap.First(a.base.Generate()))
	})
}

func (a tblArb) Shrink(p TblParams) iter.Seq[TblParams] {
	shrinks := slices.Collect(a.base.Shrink(p.V3PlainParams))
	out := make([]TblParams, 0, len(shrinks))
	for _, s := range shrinks {
		out = append(out, NewTblParams(s))
	}

	return slices.Values(out)
}
