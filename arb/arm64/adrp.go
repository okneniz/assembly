package arm64

// Generator for adrp — one generator, one type, one constructor (Adrp).

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/disasm"
)

// AdrpParams — parameters of adrp (a register and the signed 4KB page
// count from the instruction's page).
type AdrpParams struct {
	Rd   arm64.Reg
	Pages int64
}

func NewAdrpParams(rd arm64.Reg, pages int64) AdrpParams {
	return AdrpParams{
		Rd:    rd,
		Pages: pages,
	}
}

func (p AdrpParams) Instr() arm64.Instr {
	in, err := arm64.New().Adrp(p.Rd, p.Pages)
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}
func (p AdrpParams) String() string {
	return p.Instr().ObjDump(disasm.DefaultViewCtx())
}

// adrpGen — generator for adrp: an x-register and the page count uniform
// in the imm21 range (each page unit is 4KB of address space).
type adrpGen struct {
	rnd   *rand.Rand
	pages ohsnap.Arbitrary[int64]
}

func newAdrpGen(rnd *rand.Rand) adrpGen {
	return adrpGen{
		rnd:   rnd,
		pages: BrOff(rnd, 1<<20),
	}
}

// Adrp — an arbitrary adrp.
func Adrp(rnd *rand.Rand) ohsnap.Arbitrary[AdrpParams] {
	return newAdrpGen(rnd)
}

func (g adrpGen) Generate() iter.Seq[AdrpParams] {
	return arb.Stream(func() AdrpParams {
		return NewAdrpParams(genReg(g.rnd, true, false, true), ohsnap.First(g.pages.Generate()))
	})
}

func (g adrpGen) Shrink(p AdrpParams) iter.Seq[AdrpParams] {
	var out []AdrpParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewAdrpParams(r, p.Pages))
	}

	for _, v := range slices.Collect(g.pages.Shrink(p.Pages)) {
		out = append(out, NewAdrpParams(p.Rd, v))
	}

	return slices.Values(out)
}
