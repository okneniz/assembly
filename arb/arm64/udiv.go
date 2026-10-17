package arm64

// Generator for udiv — one generator, one type, one constructor (Udiv).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/arbitrary"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// UdivParams — parameters of udiv rd, rn, rm.
type UdivParams struct {
	Rd, Rn, Rm arm64.Reg
}

func NewUdivParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg) UdivParams {
	return UdivParams{
		Rd: rd,
		Rn: rn,
		Rm: rm,
	}
}

func (p UdivParams) Instr() arm64.Instr {
	in, err := arm64.New().Udiv(p.Rd, p.Rn, p.Rm)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p UdivParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// udivGen — generator for udiv: registers of the same width, 31st is zr.
type udivGen struct {
	rnd  *rand.Rand
	coin ohsnap.Arbitrary[bool]
}

// Udiv — an arbitrary udiv.
func Udiv(rnd *rand.Rand) ohsnap.Arbitrary[UdivParams] {
	return newUdivGen(rnd)
}

func newUdivGen(rnd *rand.Rand) udivGen {
	return udivGen{
		rnd:  rnd,
		coin: arbitrary.Bool(rnd),
	}
}

func (g udivGen) Generate() UdivParams {
	is64 := g.coin.Generate()
	return NewUdivParams(
		genReg(g.rnd, is64, false, true),
		genReg(g.rnd, is64, false, true),
		genReg(g.rnd, is64, false, true),
	)
}

func (g udivGen) Shrink(p UdivParams) iter.Seq[UdivParams] {
	regs := regShrunk(p.Rd)
	out := make([]UdivParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewUdivParams(r, p.Rn, p.Rm))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewUdivParams(p.Rd, r, p.Rm))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewUdivParams(p.Rd, p.Rn, r))
	}

	return slices.Values(out)
}
