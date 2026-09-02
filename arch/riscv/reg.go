package riscv

// The register operand: registers are encoded by number; the canonical
// name is the ABI name (zero/ra/sp/t0/...), as the decoder prints it.

import "fmt"

// Reg - an integer register x0..x31 (the name is an ABI alias of the number).
type Reg struct {
	num uint8
}

func newReg(num uint8) Reg {
	return Reg{num: num}
}

// X - a register by number x0..x31.
func X(n int) (Reg, error) {
	if n < 0 || n > 31 {
		return Reg{}, fmt.Errorf("riscv.X: register number %d outside 0..31", n)
	}

	return newReg(uint8(n)), nil
}

// Named ABI registers (others are available via X(n): t0 = X(5), a0 = X(10)...).
var (
	Zero = newReg(0)
	Ra   = newReg(1)
	Sp   = newReg(2)
)

// Num - the register number (rd/rs1/rs2 field: 0..31).
func (r Reg) Num() uint8 {
	return r.num
}

func (r Reg) String() string {
	return r.name()
}

// name - the canonical ABI name ("zero", "t0", "a7", ...).
func (r Reg) name() string {
	return rvRegNames[r.num]
}

// fpName - the FP-register ABI name of a register number ("ft0", "fa0",
// ...): FP operands of the constructors are taken as Reg (0..31) and
// stored/displayed as f-registers.
func fpName(r Reg) string {
	return rvFRegNames[r.num]
}
