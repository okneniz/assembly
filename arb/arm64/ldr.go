package arm64

// Generator for ldr (unsigned offset) — one generator, one type, one
// constructor (Ldr).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// LdrParams — parameters of ldr rt, [rn, #off].
type LdrParams struct {
	Rt, Rn arm64.Reg
	Off    arm64.Off
}

func NewLdrParams(rt arm64.Reg, rn arm64.Reg, off arm64.Off) LdrParams {
	return LdrParams{
		Rt:  rt,
		Rn:  rn,
		Off: off,
	}
}

func (p LdrParams) Instr() arm64.Instr {
	in, err := arm64.New().Ldr(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p LdrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// ldrGen — generator for ldr: rt is x/w (rarely zr), rn is the base x/sp;
// the offset is aligned and within the range for the width of rt (scale 8
// for x, 4 for w).
type ldrGen struct {
	rnd *rand.Rand
}

func newLdrGen(rnd *rand.Rand) ldrGen {
	return ldrGen{rnd: rnd}
}

// Ldr — an arbitrary ldr (unsigned offset).
func Ldr(rnd *rand.Rand) ohsnap.Arbitrary[LdrParams] {
	return newLdrGen(rnd)
}

func (g ldrGen) Generate() iter.Seq[LdrParams] {
	return arb.Stream(func() LdrParams {
		is64 := g.rnd.IntN(2) == 1
		align := int64(4)
		if is64 {
			align = 8
		}

		return NewLdrParams(
			genReg(g.rnd, is64, false, true),
			genReg(g.rnd, true, true, false),
			arm64.Off(g.rnd.Int64N(0x1000)*align),
		)
	})
}

func (g ldrGen) Shrink(p LdrParams) iter.Seq[LdrParams] {
	scale := uint32(2)
	if p.Rt.Is64() {
		scale = 3
	}

	var out []LdrParams
	for _, r := range regShrunk(p.Rt) {
		// changing the width of rt changes the offset scale — skip
		if r.Is64() != p.Rt.Is64() {
			continue
		}

		out = append(out, NewLdrParams(r, p.Rn, p.Off))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewLdrParams(p.Rt, r, p.Off))
	}

	for _, o := range offShrunk(p.Off, scale) {
		out = append(out, NewLdrParams(p.Rt, p.Rn, o))
	}

	return slices.Values(out)
}
