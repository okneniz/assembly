package arm64

// Generator for adds (immediate) — one generator, one type, one
// constructor (AddsImm). Unlike add, the S form reads rd 31 as zr —
// sp/wsp is not allowed there.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"
	"github.com/okneniz/oh-snap/shrink"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AddsImmParams — parameters of adds rd, rn, #imm12[, lsl #12].
type AddsImmParams struct {
	Rd, Rn arm64.Reg // 31 reads as zr
	Imm    arm64.Imm12
	Sh     arm64.Sh12
}

func NewAddsImmParams(rd arm64.Reg, rn arm64.Reg, imm arm64.Imm12, sh arm64.Sh12) AddsImmParams {
	return AddsImmParams{
		Rd:  rd,
		Rn:  rn,
		Imm: imm,
		Sh:  sh,
	}
}

func (p AddsImmParams) Instr() arm64.Instr {
	in, err := arm64.New().AddsImm(p.Rd, p.Rn, p.Imm, p.Sh)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p AddsImmParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// addsImmGen — generator for adds: same-width registers (rd carries no
// 31st), immediate 0..0xfff, shift no/lsl #12.
type addsImmGen struct {
	rnd *rand.Rand
}

// AddsImm — an arbitrary adds (immediate).
func AddsImm(rnd *rand.Rand) ohsnap.Arbitrary[AddsImmParams] {
	return newAddsImmGen(rnd)
}

func newAddsImmGen(rnd *rand.Rand) addsImmGen {
	return addsImmGen{rnd: rnd}
}

func (g addsImmGen) Generate() iter.Seq[AddsImmParams] {
	return arbStream(func() AddsImmParams {
		is64 := g.rnd.IntN(2) == 1
		return NewAddsImmParams(
			genReg(g.rnd, is64, false, true), // rd 31 reads as zr (the S form)
			genReg(g.rnd, is64, true, false), // rn 31 reads as sp
			imm12(g.rnd.Int64N(0x1000)),
			arm64.Sh12(g.rnd.IntN(2)),
		)
	})
}

func (g addsImmGen) Shrink(p AddsImmParams) iter.Seq[AddsImmParams] {
	v, err := immValue(p.Imm)
	if err != nil {
		return ohsnapEmpty[AddsImmParams]() // String() of our own type is unparseable — invariant
	}

	var out []AddsImmParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewAddsImmParams(r, p.Rn, p.Imm, p.Sh))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewAddsImmParams(p.Rd, r, p.Imm, p.Sh))
	}

	for d := range shrink.Halving[int64](0)(v) {
		imm, err := arm64.New().Imm12(d)
		if err != nil {
			continue // unreachable: half of a valid imm12 is always in 0..4095
		}

		out = append(out, NewAddsImmParams(p.Rd, p.Rn, imm, p.Sh))
	}

	if p.Sh != arm64.NoSh12 {
		out = append(out, NewAddsImmParams(p.Rd, p.Rn, p.Imm, arm64.NoSh12))
	}

	return slices.Values(out)
}
