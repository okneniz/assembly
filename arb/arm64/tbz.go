package arm64

// Generator for tbz — one generator, one type, one constructor (Tbz).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// TbzParams — parameters of tbz (a register, the tested bit and the
// pc-relative offset; bit >= 32 requires an x-register).
type TbzParams struct {
	Rt  arm64.Reg
	Bit uint32
	Off int64
}

func NewTbzParams(rt arm64.Reg, bit uint32, off int64) TbzParams {
	return TbzParams{
		Rt:  rt,
		Bit: bit,
		Off: off,
	}
}

func (p TbzParams) Instr() arm64.Instr {
	in, err := arm64.New().Tbz(p.Rt, p.Bit, p.Off)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

func (p TbzParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// tbzGen — generator for tbz: the bit chooses the register width (bit
// < 32 — w, else x), the offset uniform in the ±32KB imm14 range.
type tbzGen struct {
	rnd *rand.Rand
	off ohsnap.Arbitrary[int64]
}

// Tbz — an arbitrary tbz.
func Tbz(rnd *rand.Rand) ohsnap.Arbitrary[TbzParams] {
	return newTbzGen(rnd)
}

func newTbzGen(rnd *rand.Rand) tbzGen {
	return tbzGen{
		rnd: rnd,
		off: BrOff(rnd, 1<<15),
	}
}

func (g tbzGen) Generate() iter.Seq[TbzParams] {
	return arb.Stream(func() TbzParams {
		// the sf bit of the encoding IS the bit's b5: an x form exists
		// only for bits 32..63 (a low bit canonizes to the w form), the
		// w form covers bits 0..31.
		bit := g.rnd.IntN(32)
		if g.rnd.IntN(2) == 1 {
			bit += 32
		}

		return NewTbzParams(
			genReg(g.rnd, bit >= 32, false, true),
			uint32(bit),
			ohsnap.First(g.off.Generate()),
		)
	})
}

func (g tbzGen) Shrink(p TbzParams) iter.Seq[TbzParams] {
	var out []TbzParams
	for _, r := range regShrunk(p.Rt) {
		out = append(out, NewTbzParams(r, p.Bit, p.Off))
	}

	for bit := p.Bit / 2; ; bit /= 2 {
		if p.Rt.Is64() == (bit >= 32) {
			out = append(out, NewTbzParams(p.Rt, bit, p.Off))
		}

		if bit == 0 {
			break
		}
	}

	for _, v := range slices.Collect(g.off.Shrink(p.Off)) {
		out = append(out, NewTbzParams(p.Rt, p.Bit, v))
	}

	return slices.Values(out)
}
