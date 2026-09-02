package loong64

import (
	arch "github.com/okneniz/assembly/arch/loong64"
)

// Pre-minted ABI-named registers ($t0 = r12, $a0 = r4, ...): chain call
// sites stay arithmetic-free.

// reg - the pre-minting helper: the number is a constant 0..31, the
// constructor error is unreachable.
func reg(n int) arch.Reg {
	r, err := arch.R(n)
	if err != nil {
		panic(err) // unreachable: 0..31 is a valid register number
	}

	return r
}

var (
	Zero = reg(0)
	Ra   = reg(1)
	Tp   = reg(2)
	Sp   = reg(3)

	// $a0..$a7 - the argument registers r4..r11.
	A0 = reg(4)
	A1 = reg(5)
	A2 = reg(6)
	A3 = reg(7)
	A4 = reg(8)
	A5 = reg(9)
	A6 = reg(10)
	A7 = reg(11)

	// $t0..$t8 - the temporaries r12..r20.
	T0 = reg(12)
	T1 = reg(13)
	T2 = reg(14)
	T3 = reg(15)
	T4 = reg(16)
	T5 = reg(17)
	T6 = reg(18)
	T7 = reg(19)
	T8 = reg(20)

	Fp = reg(22)
)
