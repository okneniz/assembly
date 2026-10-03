package arm64

// Generator for cbnz — one generator, one type, one constructor (Cbnz).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// CbnzParams — parameters of cbnz (a register and the pc-relative offset).
type CbnzParams struct {
	Rt  arm64.Reg
	Off int64
}

func NewCbnzParams(rt arm64.Reg, off int64) CbnzParams {
	return CbnzParams{
		Rt:  rt,
		Off: off,
	}
}

func (p CbnzParams) Instr() arm64.Instr {
	in, err := arm64.New().Cbnz(p.Rt, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p CbnzParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// cbnzGen — generator for cbnz: an x/w register (zr allowed) and the
// offset uniform in the ±1MB imm19 range.
type cbnzGen struct {
	rnd *rand.Rand
	off ohsnap.Arbitrary[int64]
}

// Cbnz — an arbitrary cbnz.
func Cbnz(rnd *rand.Rand) ohsnap.Arbitrary[CbnzParams] {
	return newCbnzGen(rnd)
}

func newCbnzGen(rnd *rand.Rand) cbnzGen {
	return cbnzGen{
		rnd: rnd,
		off: BrOff(rnd, 1<<20),
	}
}

func (g cbnzGen) Generate() iter.Seq[CbnzParams] {
	return arb.Stream(func() CbnzParams {
		return NewCbnzParams(
			genReg(g.rnd, g.rnd.IntN(2) == 1, false, true),
			ohsnap.First(g.off.Generate()),
		)
	})
}

func (g cbnzGen) Shrink(p CbnzParams) iter.Seq[CbnzParams] {
	regs := regShrunk(p.Rt)
	out := make([]CbnzParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewCbnzParams(r, p.Off))
	}

	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewCbnzParams(p.Rt, v))
	}

	return slices.Values(out)
}
