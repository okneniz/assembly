package riscv

import (
	arch "github.com/okneniz/assembly/arch/riscv"
)

// Pre-minted ABI-named registers (zero = x0, a0 = x10, ...): chain call
// sites stay arithmetic-free.

// reg - the pre-minting helper: the number is a constant 0..31, the
// constructor error is unreachable.
func reg(n int) arch.Reg {
	r, err := arch.X(n)
	if err != nil {
		panic(err) // unreachable: 0..31 is a valid register number
	}

	return r
}

var (
	Zero = reg(0)
	Ra   = reg(1)
	Sp   = reg(2)
	Gp   = reg(3)
	Tp   = reg(4)

	// t0..t2 - the temporaries x5..x7.
	T0 = reg(5)
	T1 = reg(6)
	T2 = reg(7)

	Fp = reg(8) // s0
	S1 = reg(9)

	// a0..a7 - the argument registers x10..x17.
	A0 = reg(10)
	A1 = reg(11)
	A2 = reg(12)
	A3 = reg(13)
	A4 = reg(14)
	A5 = reg(15)
	A6 = reg(16)
	A7 = reg(17)

	// s2..s11 - the saved registers x18..x27.
	S2  = reg(18)
	S3  = reg(19)
	S4  = reg(20)
	S5  = reg(21)
	S6  = reg(22)
	S7  = reg(23)
	S8  = reg(24)
	S9  = reg(25)
	S10 = reg(26)
	S11 = reg(27)

	// t3..t6 - the temporaries x28..x31.
	T3 = reg(28)
	T4 = reg(29)
	T5 = reg(30)
	T6 = reg(31)
)
