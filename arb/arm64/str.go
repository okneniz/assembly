package arm64

// Generator for str (unsigned offset) — one generator, one type, one
// constructor (arm64.NewStr).

import (
	"math/rand/v2"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// StrParams — parameters of str rt, [rn, #off].
type StrParams struct {
	Rt, Rn arm64.Reg
	Off    arm64.Off
}

func NewStrParams(rt arm64.Reg, rn arm64.Reg, off arm64.Off) StrParams {
	return StrParams{
		Rt:  rt,
		Rn:  rn,
		Off: off,
	}
}

func (p StrParams) Instr() arm64.Instr {
	in, err := arm64.NewStr(p.Rt, p.Rn, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p StrParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// strGen — generator for str: rt is x/w (rarely zr), rn is the base x/sp;
// the offset is aligned and within the range for the width of rt (scale 8
// for x, 4 for w).
type strGen struct {
	rnd *rand.Rand
}

func newStrGen(rnd *rand.Rand) strGen {
	return strGen{rnd: rnd}
}

// Str — an arbitrary str (unsigned offset).
func Str(rnd *rand.Rand) ohsnap.Arbitrary[StrParams] {
	return newStrGen(rnd)
}

func (g strGen) Generate() StrParams {
	is64 := g.rnd.IntN(2) == 1
	align := int64(4)
	if is64 {
		align = 8
	}

	return NewStrParams(
		genReg(g.rnd, is64, false, true),
		genReg(g.rnd, true, true, false),
		arm64.Off(g.rnd.Int64N(0x1000)*align),
	)
}

func (g strGen) Shrink(p StrParams) []StrParams {
	scale := uint32(2)
	if p.Rt.Is64() {
		scale = 3
	}

	var out []StrParams
	for _, r := range regShrunk(p.Rt) {
		// changing the width of rt changes the offset scale — skip
		if r.Is64() != p.Rt.Is64() {
			continue
		}

		out = append(out, NewStrParams(r, p.Rn, p.Off))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewStrParams(p.Rt, r, p.Off))
	}

	for _, o := range offShrunk(p.Off, scale) {
		out = append(out, NewStrParams(p.Rt, p.Rn, o))
	}

	return out
}
