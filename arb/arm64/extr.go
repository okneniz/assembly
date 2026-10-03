package arm64

// Generator for extr — one generator, one type, one constructor (Extr).
// Both widths walk (the 64-bit form shifts lsb 0..63, the 32-bit one
// 0..31).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// ExtrParams — parameters of extr rd, rn, rm, #lsb.
type ExtrParams struct {
	Rd, Rn, Rm arm64.Reg // same-width registers, 31 reads as zr
	Lsb        uint32
}

func NewExtrParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg, lsb uint32) ExtrParams {
	return ExtrParams{
		Rd:  rd,
		Rn:  rn,
		Rm:  rm,
		Lsb: lsb,
	}
}

func (p ExtrParams) Instr() arm64.Instr {
	in, err := arm64.New().Extr(p.Rd, p.Rn, p.Rm, imm6(int64(p.Lsb)))
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p ExtrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// extrGen — generator for extr: registers of the same width, lsb
// 0..regsize-1.
type extrGen struct {
	rnd *rand.Rand
}

// Extr — an arbitrary extr.
func Extr(rnd *rand.Rand) ohsnap.Arbitrary[ExtrParams] {
	return newExtrGen(rnd)
}

func newExtrGen(rnd *rand.Rand) extrGen {
	return extrGen{rnd: rnd}
}

func (g extrGen) Generate() iter.Seq[ExtrParams] {
	return arbStream(func() ExtrParams {
		is64 := g.rnd.IntN(2) == 1
		lsb := uint32(g.rnd.IntN(64))
		if !is64 {
			lsb = uint32(g.rnd.IntN(32))
		}

		return NewExtrParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, is64, false, true),
			lsb,
		)
	})
}

func (g extrGen) Shrink(p ExtrParams) iter.Seq[ExtrParams] {
	regs := regShrunk(p.Rd)
	out := make([]ExtrParams, 0, len(regs))
	for _, r := range regs {
		out = append(out, NewExtrParams(r, p.Rn, p.Rm, p.Lsb))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewExtrParams(p.Rd, r, p.Rm, p.Lsb))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewExtrParams(p.Rd, p.Rn, r, p.Lsb))
	}

	for _, v := range u32Halved(p.Lsb) {
		out = append(out, NewExtrParams(p.Rd, p.Rn, p.Rm, v))
	}

	return slices.Values(out)
}
