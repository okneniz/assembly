package arm64

import (
	"fmt"
	"strconv"
)

// FReg — FP/SIMD register operand: the scalar views of the one vector
// file, s0..s31 (low 32 bits) and d0..d31 (low 64 bits). Unlike the
// integer file there is no named 31st register — v31/s31/d31 is an
// ordinary register.
type FReg struct {
	num  uint8 // 0..31
	is64 bool  // d = true, s = false
}

// newFReg - the register constructor: the struct is assembled only here.
func newFReg(num uint8, is64 bool) FReg {
	return FReg{
		num:  num,
		is64: is64,
	}
}

// S — single-precision register s0..s31.
func S(n int) (FReg, error) {
	if n < 0 || n > 31 {
		return FReg{}, fmt.Errorf(
			"arm64.S: register number %d is out of 0..31", n,
		)
	}

	return newFReg(uint8(n), false), nil
}

// D — double-precision register d0..d31.
func D(n int) (FReg, error) {
	if n < 0 || n > 31 {
		return FReg{}, fmt.Errorf(
			"arm64.D: register number %d is out of 0..31", n,
		)
	}

	return newFReg(uint8(n), true), nil
}

// FRegOf — an FP register by its source name (s0, d31): the inverse of
// (FReg).name.
func FRegOf(name string) (FReg, error) {
	if len(name) < 2 {
		return FReg{}, fmt.Errorf("arm64.FRegOf: unknown register %q", name)
	}

	n, err := strconv.Atoi(name[1:])
	if err != nil || n < 0 || n > 31 {
		return FReg{}, fmt.Errorf("arm64.FRegOf: unknown register %q", name)
	}

	switch name[0] {
	case 's':
		return S(n)
	case 'd':
		return D(n)
	}

	return FReg{}, fmt.Errorf("arm64.FRegOf: unknown register %q", name)
}

// Num — the register number (0..31).
func (r FReg) Num() uint8 {
	return r.num
}

// Is64 — width of the view: d = true, s = false.
func (r FReg) Is64() bool {
	return r.is64
}

// kind — the fpKind of the operand.
func (r FReg) kind() fpKind {
	if r.is64 {
		return kD
	}

	return kS
}

func (r FReg) String() string {
	return r.name()
}

// name - canonical name ("s5", "d7").
func (r FReg) name() string {
	if r.is64 {
		return "d" + strconv.Itoa(int(r.num))
	}

	return "s" + strconv.Itoa(int(r.num))
}
