package arm64

// Generator for cbz — one generator, one type, one constructor (Cbz).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// CbzParams — parameters of cbz (a register and the pc-relative offset).
type CbzParams struct {
	Rt  arm64.Reg
	Off int64
}

func NewCbzParams(rt arm64.Reg, off int64) CbzParams {
	return CbzParams{
		Rt:  rt,
		Off: off,
	}
}

func (p CbzParams) Instr() arm64.Instr {
	in, err := arm64.New().Cbz(p.Rt, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p CbzParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// cbzGen — generator for cbz: an x/w register (zr allowed) and the offset
// uniform in the ±1MB imm19 range.
type cbzGen struct {
	rnd *rand.Rand
	off ohsnap.Arbitrary[int64]
}

func newCbzGen(rnd *rand.Rand) cbzGen {
	return cbzGen{
		rnd: rnd,
		off: BrOff(rnd, 1<<20),
	}
}

// Cbz — an arbitrary cbz.
func Cbz(rnd *rand.Rand) ohsnap.Arbitrary[CbzParams] {
	return newCbzGen(rnd)
}

func (g cbzGen) Generate() iter.Seq[CbzParams] {
	return arb.Stream(func() CbzParams {
		return NewCbzParams(genReg(g.rnd, g.rnd.IntN(2) == 1, false, true), ohsnap.First(g.off.Generate()))
	})
}

func (g cbzGen) Shrink(p CbzParams) iter.Seq[CbzParams] {
	var out []CbzParams
	for _, r := range regShrunk(p.Rt) {
		out = append(out, NewCbzParams(r, p.Off))
	}

	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewCbzParams(p.Rt, v))
	}

	return slices.Values(out)
}
