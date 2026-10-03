package arm64

// Generator for rev32 — one generator, one type, one constructor (Rev32).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// Rev32Params — parameters of rev32 rd, rn (the 64-bit form only).
type Rev32Params struct {
	Rd, Rn arm64.Reg
}

func NewRev32Params(rd arm64.Reg, rn arm64.Reg) Rev32Params {
	return Rev32Params{
		Rd: rd,
		Rn: rn,
	}
}

func (p Rev32Params) Instr() arm64.Instr {
	in, err := arm64.New().Rev32(p.Rd, p.Rn)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p Rev32Params) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// rev32Gen — generator for rev32: x-registers, occasionally xzr.
type rev32Gen struct {
	rnd *rand.Rand
}

// Rev32 — an arbitrary rev32.
func Rev32(rnd *rand.Rand) ohsnap.Arbitrary[Rev32Params] {
	return newRev32Gen(rnd)
}

func newRev32Gen(rnd *rand.Rand) rev32Gen {
	return rev32Gen{rnd: rnd}
}

func (g rev32Gen) Generate() iter.Seq[Rev32Params] {
	return arb.Stream(func() Rev32Params {
		return NewRev32Params(genReg(g.rnd, true, false, true), genReg(g.rnd, true, false, true))
	})
}

func (g rev32Gen) Shrink(p Rev32Params) iter.Seq[Rev32Params] {
	regs := regShrunk(p.Rd)
	out := make([]Rev32Params, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewRev32Params(r, p.Rn))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewRev32Params(p.Rd, r))
	}

	return slices.Values(out)
}
