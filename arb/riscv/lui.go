package riscv

// Generator for lui — one generator, one type, one constructor
// (Lui).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// LuiParams — parameters of lui.
type LuiParams struct {
	Rd  riscv.Reg
	Imm riscv.Imm20
}

func NewLuiParams(rd riscv.Reg, imm riscv.Imm20) LuiParams {
	return LuiParams{
		Rd:  rd,
		Imm: imm,
	}
}

func (p LuiParams) Instr() riscv.Instr {
	return riscv.New().Lui(p.Rd, p.Imm)
}

func (p LuiParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// luiGen — generator for lui: register x0..x31, U-type field 0..0xfffff.
type luiGen struct {
	rnd *rand.Rand
}

// Lui — an arbitrary lui.
func Lui(rnd *rand.Rand) ohsnap.Arbitrary[LuiParams] {
	return newLuiGen(rnd)
}

func newLuiGen(rnd *rand.Rand) luiGen {
	return luiGen{rnd: rnd}
}

func (g luiGen) Generate() iter.Seq[LuiParams] {
	return arb.Stream(func() LuiParams {
		// the c.lui compression window (imm20 the 6-bit form reaches:
		// 1..31 and the sext-negative tail) is sampled head-on - the
		// uniform 0..0xfffff draw almost never lands there
		if g.rnd.IntN(4) == 0 {
			imm := uint32(1 + g.rnd.IntN(31))
			if g.rnd.IntN(2) == 1 {
				imm = 0x100000 - imm // 0xfffff..0xffff0
			}

			if v, err := riscv.New().Imm20(int64(imm)); err == nil {
				return NewLuiParams(reg(g.rnd), v)
			}
		}

		return NewLuiParams(reg(g.rnd), imm20(g.rnd))
	})
}

func (g luiGen) Shrink(p LuiParams) iter.Seq[LuiParams] {
	rd := regShrunk(p.Rd)
	imms := immShrunk(p.Imm, riscv.New().Imm20, imm20Shrink)
	out := make([]LuiParams, 0, len(rd)+len(imms))
	for _, r := range rd {
		out = append(out, NewLuiParams(r, p.Imm))
	}

	for _, v := range imms {
		out = append(out, NewLuiParams(p.Rd, v))
	}

	return slices.Values(out)
}
