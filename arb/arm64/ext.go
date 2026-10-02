package arm64

// The extended-register (add/sub … rm, ext #imm3) family core: one
// parameter shape, one generator; the four mnemonics (add/adds/sub/subs)
// are thin constructors over it. The extension names follow the register
// width — a w form has no uxtx/sxtx (the ctor does not check the pairing,
// the generator does).

import (
	"github.com/okneniz/assembly/arch/arm64"
	"iter"
	"math/rand/v2"
	"slices"

	"github.com/okneniz/oh-snap/shrink"
)

// ExtParams — parameters of add rd, rn, rm, ext #imm3 (and its
// adds/sub/subs twins).
type ExtParams struct {
	Rd, Rn, Rm arm64.Reg // 31 reads as sp/wsp
	Ext        string
	Imm3       uint32
}

func NewExtParams(rd arm64.Reg, rn arm64.Reg, rm arm64.Reg, ext string, imm3 uint32) ExtParams {
	return ExtParams{
		Rd:   rd,
		Rn:   rn,
		Rm:   rm,
		Ext:  ext,
		Imm3: imm3,
	}
}

// extNames — the extension names of the register width.
func extNames(is64 bool) []string {
	if is64 {
		return []string{"uxtb", "uxth", "uxtw", "uxtx", "sxtb", "sxth", "sxtw", "sxtx"}
	}

	return []string{"uxtb", "uxth", "uxtw", "sxtb", "sxth", "sxtw"}
}

// extGen — the family generator: registers of one width, a
// width-compatible extension, imm3 0..7. The S forms (adds/subs) read
// rd 31 as zr — the generator keeps the plain register there.
type extGen struct {
	rnd   *rand.Rand
	sForm bool
}

func newExtGen(rnd *rand.Rand, sForm bool) extGen {
	return extGen{
		rnd:   rnd,
		sForm: sForm,
	}
}

// ext — the shared Generate/Shrink core of add/sub ext.
func ext(rnd *rand.Rand) extGen {
	return newExtGen(rnd, false)
}

// extS — the core of adds/subs ext (rd 31 reads as zr).
func extS(rnd *rand.Rand) extGen {
	return newExtGen(rnd, true)
}

func (g extGen) Generate() iter.Seq[ExtParams] {
	return arbStream(func() ExtParams {
		is64 := g.rnd.IntN(2) == 1
		names := extNames(is64)
		ext := names[g.rnd.IntN(len(names))]

		// rd reads as sp only in the plain add/sub x forms with a
		// 64-bit extension (the sp arithmetic rule); rm is x/w with zr
		// (never sp - clang refuses it there)
		rdSp := !g.sForm && is64 && (ext == "uxtx" || ext == "sxtx")
		return NewExtParams(
			genReg(g.rnd, is64, rdSp, false),
			genReg(g.rnd, is64, true, false),
			genReg(g.rnd, is64, false, true),
			ext,
			uint32(g.rnd.IntN(8)),
		)
	})
}

func (g extGen) Shrink(p ExtParams) iter.Seq[ExtParams] {
	var out []ExtParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewExtParams(r, p.Rn, p.Rm, p.Ext, p.Imm3))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewExtParams(p.Rd, r, p.Rm, p.Ext, p.Imm3))
	}

	for _, r := range regShrunk(p.Rm) {
		out = append(out, NewExtParams(p.Rd, p.Rn, r, p.Ext, p.Imm3))
	}

	for d := range shrink.Halving[uint32](0)(p.Imm3) {
		out = append(out, NewExtParams(p.Rd, p.Rn, p.Rm, p.Ext, d))
	}

	return slices.Values(out)
}
