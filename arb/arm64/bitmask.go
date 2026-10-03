package arm64

// The logical immediate family core (and/orr/eor/ands rd, rn, #imm): the
// structural generation of the bitmask — a run of Len ones rotated by Rot
// inside the element of size Esz, replicated over the register width.
// Every (esize, len < esize, rot < esize) triple is encodable; the
// all-ones value (len == esize) is not a logical immediate and never
// leaves the generator. The ctor takes the VALUE, so the params keep the
// structural axes — the shrink moves them, not the opaque number.

import (
	"iter"
	"math/rand/v2"
	"slices"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arch/arm64"
)

// BitmaskParams — the structural axes of a logical immediate.
type BitmaskParams struct {
	Rd, Rn arm64.Reg // 31 reads as zr
	Esz    uint32    // the element size, a power of two (2..64; 2..32 for w forms)
	Len    uint32    // the run of ones, 1..Esz-1
	Rot    uint32    // the rotation of the run, 0..Esz-1
}

func NewBitmaskParams(
	rd arm64.Reg,
	rn arm64.Reg,
	esz uint32,
	length uint32,
	rot uint32,
) BitmaskParams {
	return BitmaskParams{
		Rd:  rd,
		Rn:  rn,
		Esz: esz,
		Len: length,
		Rot: rot,
	}
}

// Value — the immediate the axes encode.
func (p BitmaskParams) Value() uint64 {
	mask := uint64(1)<<p.Esz - 1
	if p.Esz == 64 {
		mask = ^uint64(0)
	}

	ones := (uint64(1)<<p.Len - 1) & mask
	r := p.Rot % p.Esz
	rotated := ((ones << r) | (ones >> (p.Esz - r))) & mask
	v := rotated
	for w := p.Esz; w < 64; w *= 2 {
		v |= v << w
	}

	return v
}

// bitmaskGen — the family generator.
type bitmaskGen struct {
	rnd *rand.Rand
}

func (g bitmaskGen) Generate() iter.Seq[BitmaskParams] {
	return arbStream(func() BitmaskParams {
		is64 := g.rnd.IntN(2) == 1
		choices := eszChoices(is64)
		esz := choices[g.rnd.IntN(len(choices))]
		return NewBitmaskParams(
			genReg(
				g.rnd,
				is64,
				false,
				false,
			), // clang refuses a zr destination (the encoding is fine)
			genReg(g.rnd, is64, false, true),
			esz,
			uint32(g.rnd.IntN(int(esz)-1))+1, // 1..esz-1
			uint32(g.rnd.IntN(int(esz))),
		)
	})
}

func (g bitmaskGen) Shrink(p BitmaskParams) iter.Seq[BitmaskParams] {
	var out []BitmaskParams
	for _, r := range regShrunk(p.Rd) {
		out = append(out, NewBitmaskParams(r, p.Rn, p.Esz, p.Len, p.Rot))
	}

	for _, r := range regShrunk(p.Rn) {
		out = append(out, NewBitmaskParams(p.Rd, r, p.Esz, p.Len, p.Rot))
	}

	// a smaller element that still fits the run (the rotation folds)
	for esz := p.Esz / 2; esz >= 2 && esz > p.Len; esz /= 2 {
		out = append(out, NewBitmaskParams(p.Rd, p.Rn, esz, p.Len, p.Rot%esz))
	}

	for _, v := range u32Halved(p.Len) {
		if v == 0 {
			continue // the run keeps >= 1
		}

		out = append(out, NewBitmaskParams(p.Rd, p.Rn, p.Esz, v, p.Rot))
	}

	for _, v := range u32Halved(p.Rot) {
		out = append(out, NewBitmaskParams(p.Rd, p.Rn, p.Esz, p.Len, v))
	}

	return slices.Values(out)
}

// Bitmask — the exported face of the structural generator (the alias
// package builds its immediate forms on top of the same axes).
func Bitmask(rnd *rand.Rand) ohsnap.Arbitrary[BitmaskParams] {
	return newBitmaskGen(rnd)
}

func newBitmaskGen(rnd *rand.Rand) bitmaskGen {
	return bitmaskGen{rnd: rnd}
}

// bitmask — the shared Generate/Shrink core of the four logical-imm families.
func bitmask(rnd *rand.Rand) bitmaskGen {
	return newBitmaskGen(rnd)
}

// eszChoices — the element sizes of the register width, ascending.
func eszChoices(is64 bool) []uint32 {
	if is64 {
		return []uint32{2, 4, 8, 16, 32, 64}
	}

	return []uint32{2, 4, 8, 16, 32}
}

// u32Halved — the halving-toward-zero candidates of a uint32 axis.
func u32Halved(v uint32) []uint32 {
	if v == 0 {
		return nil
	}

	var out []uint32
	for d := v / 2; ; d /= 2 {
		out = append(out, d)
		if d == 0 {
			break
		}
	}

	return out
}
