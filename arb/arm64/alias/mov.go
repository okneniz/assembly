package alias

// Generator for the mov alias — one generator, one type, one text form
// family: mov rd, rm | mov rd, #imm (the three immediate classes the
// alias layer accepts: a movz lane, a negative movn value, and the
// all-ones-except-one-lane positive movn pattern).

import (
	"iter"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"

	ohsnap "github.com/okneniz/oh-snap"

	a64 "github.com/okneniz/assembly/arb/arm64"
	arm64 "github.com/okneniz/assembly/arch/arm64"
)

// The immediate classes of mov (their encoding is movz/movn).
const (
	movClassZ  = 0 // v = lane << 16hw (movz)
	movClassN  = 1 // v = -(lane << 16hw) - 1 (movn, negative)
	movClassPN = 2 // v = ^(lane << 16hw) (movn, all ones except one lane)
)

// MovParams — parameters of the mov alias.
type MovParams struct {
	Rd, Rm arm64.Reg
	Lane   int64
	Hw     uint32
	Class  int
	IsReg  bool
}

func NewMovParams(rd arm64.Reg, rm arm64.Reg, lane int64, hw uint32, class int, isReg bool) MovParams {
	return MovParams{
		Rd:    rd,
		Rm:    rm,
		Lane:  lane,
		Hw:    hw,
		Class: class,
		IsReg: isReg,
	}
}

func (p MovParams) String() string {
	if p.IsReg {
		return "mov " + p.Rd.String() + ", " + p.Rm.String()
	}

	lane := uint64(p.Lane) << (16 * p.Hw)
	switch p.Class {
	case movClassN:
		return "mov " + p.Rd.String() + ", #" + strconv.FormatInt(-int64(lane)-1, 10)
	case movClassPN:
		return "mov " + p.Rd.String() + ", " + movImmText(^lane)
	default:
		return "mov " + p.Rd.String() + ", " + movImmText(lane)
	}
}

// movImmText — the immediate text: decimal while it fits int64, hex
// beyond (the decimal lexer stops at int64; objdump prints these hex).
func movImmText(v uint64) string {
	if v > math.MaxInt64 {
		return "#0x" + strconv.FormatUint(v, 16)
	}

	return "#" + strconv.FormatUint(v, 10)
}

func (p MovParams) Instr() arm64.Instr {
	in, err := instrOfText(p.String())
	if err != nil {
		return nil // unreachable: fields are produced by a valid generator/shrink
	}

	return in
}

// movGen — generator for mov: the register form and the three immediate
// classes of the width (the 32-bit form knows only hw 0/1).
type movGen struct {
	rnd *rand.Rand
}

func newMovGen(rnd *rand.Rand) movGen {
	return movGen{rnd: rnd}
}

// Mov — an arbitrary mov.
func Mov(rnd *rand.Rand) ohsnap.Arbitrary[MovParams] {
	return newMovGen(rnd)
}

func (g movGen) Generate() iter.Seq[MovParams] {
	return stream(func() MovParams {
		is64 := g.rnd.IntN(2) == 1
		if g.rnd.IntN(2) == 1 {
			return NewMovParams(
				a64.GenReg(g.rnd, is64, false, true),
				a64.GenReg(g.rnd, is64, false, true),
				0,
				0,
				movClassZ,
				true,
			)
		}

		maxHw := uint32(1)
		if is64 {
			maxHw = 3
		}

		class := g.rnd.IntN(2)
		if is64 && g.rnd.IntN(2) == 1 {
			class = movClassPN // the all-ones pattern is a 64-bit form
		}

		lane := g.rnd.Int64N(0x10000)
		if class == movClassPN && lane == 0 {
			lane = 1 // the all-ones pattern needs a nonzero lane
		}

		return NewMovParams(
			a64.GenReg(g.rnd, is64, false, true),
			a64.GenReg(g.rnd, is64, false, true),
			lane,
			uint32(g.rnd.IntN(int(maxHw) + 1)),
			class,
			false,
		)
	})
}

func (g movGen) Shrink(p MovParams) iter.Seq[MovParams] {
	var out []MovParams
	for _, r := range a64.RegShrunk(p.Rd) {
		out = append(out, NewMovParams(r, p.Rm, p.Lane, p.Hw, p.Class, p.IsReg))
	}

	for _, r := range a64.RegShrunk(p.Rm) {
		out = append(out, NewMovParams(p.Rd, r, p.Lane, p.Hw, p.Class, p.IsReg))
	}

	for _, v := range halved(p.Lane) {
		if p.Class == movClassPN && v == 0 {
			continue // the all-ones pattern keeps a nonzero lane
		}

		out = append(out, NewMovParams(p.Rd, p.Rm, v, p.Hw, p.Class, p.IsReg))
	}

	if p.Hw != 0 {
		out = append(out, NewMovParams(p.Rd, p.Rm, p.Lane, 0, p.Class, p.IsReg))
	}

	if p.Class != movClassZ {
		out = append(out, NewMovParams(p.Rd, p.Rm, p.Lane, 0, movClassZ, p.IsReg))
	}

	if p.IsReg {
		out = append(out, NewMovParams(p.Rd, p.Rd, 0, 0, movClassZ, false))
	}

	return slices.Values(out)
}
