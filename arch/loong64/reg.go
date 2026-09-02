package loong64

// The register operand: instructions store register NUMBERS (a byte
// each, not names); the canonical $-name is only produced at print time.

import "fmt"

// Reg - an integer register $r0..$r31 (the name is an ABI alias of the number).
type Reg struct {
	num uint8
}

func newReg(num uint8) Reg {
	return Reg{
		num: num,
	}
}

// R - a register by number $r0..$r31.
func R(n int) (Reg, error) {
	if n < 0 || n > 31 {
		return Reg{}, fmt.Errorf("loong64.R: register number %d outside 0..31", n)
	}

	return newReg(uint8(n)), nil
}

// Named ABI registers (others are available via R(n): t0 = R(12), a0 = R(4)...).
var (
	Zero = newReg(0)
	Ra   = newReg(1)
	Tp   = newReg(2)
	Sp   = newReg(3)
	Fp   = newReg(22)
)

// Num - the register number (rd/rj/rk field: 0..31).
func (r Reg) Num() uint8 {
	return r.num
}

// String - the printed form ("$zero", "$a0", "$r21", ...).
func (r Reg) String() string {
	return laRegName(r.num)
}
