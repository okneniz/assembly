package arm64

// Generator for adr — one generator, one type, one constructor (Adr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AdrParams — parameters of adr (a register and the signed byte offset
// from the instruction's own address).
type AdrParams struct {
	Rd  arm64.Reg
	Off int64
}

func NewAdrParams(rd arm64.Reg, off int64) AdrParams {
	return AdrParams{
		Rd:  rd,
		Off: off,
	}
}

func (p AdrParams) Instr() arm64.Instr {
	in, err := arm64.New().Adr(p.Rd, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p AdrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// adrGen — generator for adr: an x-register and the offset uniform in
// the imm21 byte range (the shared branch-offset grid — 4-aligned — is
// the valid subset; the unaligned bytes ride the unit tests of Adr).
type adrGen struct {
	rnd *rand.Rand
	off ohsnap.Arbitrary[int64]
}

// Adr — an arbitrary adr.
func Adr(rnd *rand.Rand) ohsnap.Arbitrary[AdrParams] {
	return newAdrGen(rnd)
}

func newAdrGen(rnd *rand.Rand) adrGen {
	return adrGen{
		rnd: rnd,
		off: BrOff(rnd, 1<<20),
	}
}

func (g adrGen) Generate() iter.Seq[AdrParams] {
	return arb.Stream(func() AdrParams {
		return NewAdrParams(genReg(g.rnd, true, false, true), ohsnap.First(g.off.Generate()))
	})
}

func (g adrGen) Shrink(p AdrParams) iter.Seq[AdrParams] {
	regs := regShrunk(p.Rd)
	out := make([]AdrParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewAdrParams(r, p.Off))
	}

	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewAdrParams(p.Rd, v))
	}

	return slices.Values(out)
}
