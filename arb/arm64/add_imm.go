// Package arm64 — arbitrary test generators (oh-snap) for ARM64 instructions:
// one generator, parameter type, and shrink per instruction.
package arm64

// Generator for add (immediate) — one generator, one type, one constructor
// (arm64.NewAddImm).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddImmParams — parameters of add rd, rn, #imm12[, lsl #12].
type AddImmParams struct {
	Rd, Rn arm64.Reg
	Imm    arm64.Imm12
	Sh     arm64.Sh12
}

func NewAddImmParams(rd arm64.Reg, rn arm64.Reg, imm arm64.Imm12, sh arm64.Sh12) AddImmParams {
	return AddImmParams{
		Rd:  rd,
		Rn:  rn,
		Imm: imm,
		Sh:  sh,
	}
}

func (p AddImmParams) Instr() arm64.Instr {
	in, err := arm64.NewAddImm(p.Rd, p.Rn, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AddImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// addImmGen — generator for add: registers of the same width (31st = sp/wsp),
// immediate 0..0xfff, shift no/lsl#12.
type addImmGen struct {
	rnd *rand.Rand
}

func newAddImmGen(rnd *rand.Rand) addImmGen {
	return addImmGen{rnd: rnd}
}

// AddImm — an arbitrary add (immediate).
func AddImm(rnd *rand.Rand) ohsnap.Arbitrary[AddImmParams] {
	return newAddImmGen(rnd)
}

func (g addImmGen) Generate() AddImmParams {
	is64 := g.rnd.IntN(2) == 1
	return NewAddImmParams(
		genReg(g.rnd, is64, true, false),
		genReg(g.rnd, is64, true, false),
		imm12(g.rnd.Int64N(0x1000)),
		arm64.Sh12(g.rnd.IntN(2)),
	)
}

func (g addImmGen) Shrink(p AddImmParams) []AddImmParams {
	v, err := immValue(p.Imm)
	if err != nil {
		return nil // String() of our own type is unparseable — invariant
	}

	var out []AddImmParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewAddImmParams(r, p.Rn, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewAddImmParams(p.Rd, r, p.Imm, p.Sh))
	}

	for _, d := range arb.Halved(v) {
		imm, err := arm64.NewImm12(d)
		if err != nil {
			continue // unreachable: half of a valid imm12 is always in 0..4095
		}

		out = append(out, NewAddImmParams(p.Rd, p.Rn, imm, p.Sh))
	}

	if p.Sh != arm64.NoSh12 {
		out = append(out, NewAddImmParams(p.Rd, p.Rn, p.Imm, arm64.NoSh12))
	}

	return out
}
