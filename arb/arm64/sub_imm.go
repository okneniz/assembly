package arm64

// Generator for sub (immediate) — one generator, one type, one constructor
// (arm64.NewSubImm).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// SubImmParams — parameters of sub rd, rn, #imm12[, lsl #12].
type SubImmParams struct {
	Rd, Rn arm64.Reg
	Imm    arm64.Imm12
	Sh     arm64.Sh12
}

func NewSubImmParams(rd arm64.Reg, rn arm64.Reg, imm arm64.Imm12, sh arm64.Sh12) SubImmParams {
	return SubImmParams{
		Rd:  rd,
		Rn:  rn,
		Imm: imm,
		Sh:  sh,
	}
}

func (p SubImmParams) Instr() arm64.Instr {
	in, err := arm64.NewSubImm(p.Rd, p.Rn, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p SubImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// subImmGen — generator for sub: registers of the same width (31st = sp/wsp),
// immediate 0..0xfff, shift no/lsl#12.
type subImmGen struct {
	rnd *rand.Rand
}

func newSubImmGen(rnd *rand.Rand) subImmGen {
	return subImmGen{rnd: rnd}
}

// SubImm — an arbitrary sub (immediate).
func SubImm(rnd *rand.Rand) ohsnap.Arbitrary[SubImmParams] {
	return newSubImmGen(rnd)
}

func (g subImmGen) Generate() SubImmParams {
	is64 := g.rnd.IntN(2) == 1
	return NewSubImmParams(
		genReg(g.rnd, is64, true, false),
		genReg(g.rnd, is64, true, false),
		imm12(g.rnd.Int64N(0x1000)),
		arm64.Sh12(g.rnd.IntN(2)),
	)
}

func (g subImmGen) Shrink(p SubImmParams) []SubImmParams {
	v, err := immValue(p.Imm)
	if err != nil {
		return nil // String() of our own type is unparseable — invariant
	}

	var out []SubImmParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewSubImmParams(r, p.Rn, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewSubImmParams(p.Rd, r, p.Imm, p.Sh))
	}

	for _, d := range arb.Halved(v) {
		imm, err := arm64.NewImm12(d)
		if err != nil {
			continue // unreachable: half of a valid imm12 is always in 0..4095
		}

		out = append(out, NewSubImmParams(p.Rd, p.Rn, imm, p.Sh))
	}

	if p.Sh != arm64.NoSh12 {
		out = append(out, NewSubImmParams(p.Rd, p.Rn, p.Imm, arm64.NoSh12))
	}

	return out
}
