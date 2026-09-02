package arm64

// Register generators: an arbitrary register of any width and class
// (including the named 31st ones — sp/zr, as allowed in the position);
// positions with instruction-specific constraints are filtered by the
// instruction generator. The shrink is semantic — toward x0/w0 of the
// same width, then the halved number of the same class.

import (
	"iter"
	"math/rand/v2"
	"slices"
	"strconv"

	ohsnap "github.com/okneniz/oh-snap"

	"github.com/okneniz/assembly/arb"
	"github.com/okneniz/assembly/arch/arm64"
)

// regArb — an arbitrary register of any width and class.
type regArb struct {
	rnd *rand.Rand
}

func newRegArb(rnd *rand.Rand) regArb {
	return regArb{rnd: rnd}
}

// Reg — an arbitrary register of any width and class (including the named
// 31st ones); positions with instruction-specific constraints are filtered
// by the instruction generator.
func Reg(rnd *rand.Rand) ohsnap.Arbitrary[arm64.Reg] {
	return newRegArb(rnd)
}

func (a regArb) Generate() iter.Seq[arm64.Reg] {
	return arb.Stream(func() arm64.Reg {
		return genReg(a.rnd, a.rnd.IntN(2) == 1, true, true)
	})
}

func (a regArb) Shrink(r arm64.Reg) iter.Seq[arm64.Reg] {
	return slices.Values(regShrunk(r))
}

// xreg/wreg — x#/w# by number; the unreachable error branch — the input is
// bounded by construction (rnd.IntN(31), n/2 of a valid register number).
func xreg(n int) arm64.Reg {
	r, err := arm64.X(n)
	if err != nil {
		return arm64.XZR // unreachable: n is always in 0..30
	}

	return r
}

func wreg(n int) arm64.Reg {
	r, err := arm64.W(n)
	if err != nil {
		return arm64.WZR // unreachable: n is always in 0..30
	}

	return r
}

// genReg — a register of the given width: an x/w-number, occasionally the
// named 31st (sp or zr — as allowed in the position). Dependent part of
// generation: the width is decided by the instruction, not the operand generator.
func genReg(rnd *rand.Rand, is64 bool, spOK, zrOK bool) arm64.Reg {
	if spOK && rnd.IntN(16) == 0 {
		if is64 {
			return arm64.SP
		}

		return arm64.WSP
	}

	if zrOK && rnd.IntN(16) == 0 {
		if is64 {
			return arm64.XZR
		}

		return arm64.WZR
	}

	n := rnd.IntN(31)
	if is64 {
		return xreg(n)
	}

	return wreg(n)
}

// regShrunk — register shrink candidates: x0/w0 of the same width and the
// halved number of the same class. The width is preserved — the invariant
// "instruction registers are of the same width" is not violated.
func regShrunk(r arm64.Reg) []arm64.Reg {
	zero := wreg(0)
	if r.Is64() {
		zero = xreg(0)
	}

	var out []arm64.Reg
	if r != zero {
		out = append(out, zero)
	}

	if n, ok := plainRegNum(r); ok && n > 1 {
		var half arm64.Reg
		if r.Is64() {
			half = xreg(n / 2)
		} else {
			half = wreg(n / 2)
		}

		out = append(out, half)
	}

	return out
}

// plainRegNum — the number of a plain register (x#/w#); not for named ones.
func plainRegNum(r arm64.Reg) (int, bool) {
	s := r.String()
	if len(s) < 2 || (s[0] != 'x' && s[0] != 'w') {
		return 0, false
	}

	n, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, false
	}

	return n, true
}
