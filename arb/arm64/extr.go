package arm64

// Generator for extr — one generator, one type, one constructor (Extr).
// Only the 64-bit form exists in the ctor (the ror alias is the w spell);
// lsb stays below 64.

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
	Rd, Rn, Rm arm64.Reg // x-registers, 31 reads as zr
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

// extrGen — generator for extr: x-registers, lsb 0..63.
type extrGen struct {
	rnd *rand.Rand
}

func newExtrGen(rnd *rand.Rand) extrGen {
	return extrGen{rnd: rnd}
}

// Extr — an arbitrary extr.
func Extr(rnd *rand.Rand) ohsnap.Arbitrary[ExtrParams] {
	return newExtrGen(rnd)
}

func (g extrGen) Generate() iter.Seq[ExtrParams] {
	return arbStream(func() ExtrParams {
		return NewExtrParams(
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
			genReg(g.rnd, true, false, true),
			uint32(g.rnd.IntN(64)),
		)
	})
}

func (g extrGen) Shrink(p ExtrParams) iter.Seq[ExtrParams] {
	var out []ExtrParams
	for _, r := range regShrunk(p.Rd) {
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
