package riscv

// Generator for fence - one generator, one type, one constructor
// (Fence): the plain canon (fm=0 - ObjDump prints bare fence).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/disasm"
)

// FenceParams — parameters of fence.
type FenceParams struct {
	Fm uint8
}

func NewFenceParams(fm uint8) FenceParams {
	return FenceParams{Fm: fm}
}

func (p FenceParams) Instr() riscv.Instr {
	return riscv.New().Fence(p.Fm)
}

func (p FenceParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// Fence — an arbitrary fence.
func Fence(rnd *rand.Rand) ohsnap.Arbitrary[FenceParams] {
	return fenceGen{rnd: rnd}
}

type fenceGen struct {
	rnd *rand.Rand
}

func (g fenceGen) Generate() iter.Seq[FenceParams] {
	return arb.Stream(func() FenceParams {
		return NewFenceParams(0)
	})
}

func (g fenceGen) Shrink(FenceParams) iter.Seq[FenceParams] {
	var out []FenceParams
	return slices.Values(out)
}
