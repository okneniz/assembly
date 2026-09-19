package arm64

import (
	arch "github.com/okneniz/assembly/arch/arm64"
)

// Pre-minted registers: chain call sites stay arithmetic-free. The named
// 31st registers come from the arch package directly.

// reg - the pre-minting helper: the number is a constant 0..30, the
// constructor error is unreachable.
func reg(n int) arch.Reg {
	r, err := arch.X(n)
	if err != nil {
		panic(err) // unreachable: 0..30 is a valid register number
	}

	return r
}

var (
	X0 = reg(0)
	X1 = reg(1)
	X2 = reg(2)
	X3 = reg(3)
	X4 = reg(4)
	X5 = reg(5)
	X6 = reg(6)
	X7 = reg(7)

	X16 = reg(16)
	X29 = reg(29)
	X30 = reg(30)

	// XZR/SP/WSP - the named 31st registers, re-exported for one import.
	XZR = arch.XZR
	SP  = arch.SP
)

// dreg/sreg - the FP pre-minting helpers: the number is a constant
// 0..7 (the AAPCS64 argument registers), the constructor error is
// unreachable.
func dreg(n int) arch.FReg {
	r, err := arch.D(n)
	if err != nil {
		panic(err) // unreachable: 0..7 is a valid register number
	}

	return r
}

func sreg(n int) arch.FReg {
	r, err := arch.S(n)
	if err != nil {
		panic(err) // unreachable: 0..7 is a valid register number
	}

	return r
}

// vreg - the vector pre-minting helper: the number is a constant
// 0..7, the constructor error is unreachable.
func vreg(n int) arch.VReg {
	r, err := arch.V(n)
	if err != nil {
		panic(err) // unreachable: 0..7 is a valid register number
	}

	return r
}

// V0-V7 - the vector argument registers (the SIMD halves of the
// AAPCS64 argument file), the counterparts of X0-X7/D0-D7.
var (
	V0 = vreg(0)
	V1 = vreg(1)
	V2 = vreg(2)
	V3 = vreg(3)
	V4 = vreg(4)
	V5 = vreg(5)
	V6 = vreg(6)
	V7 = vreg(7)
)

// D0-D7/S0-S7 - the FP argument registers (v0-v7), the AAPCS64
// counterparts of X0-X7.
var (
	D0 = dreg(0)
	D1 = dreg(1)
	D2 = dreg(2)
	D3 = dreg(3)
	D4 = dreg(4)
	D5 = dreg(5)
	D6 = dreg(6)
	D7 = dreg(7)

	S0 = sreg(0)
	S1 = sreg(1)
	S2 = sreg(2)
	S3 = sreg(3)
	S4 = sreg(4)
	S5 = sreg(5)
	S6 = sreg(6)
	S7 = sreg(7)
)
