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
