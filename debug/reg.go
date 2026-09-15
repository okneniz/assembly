package debug

import "fmt"

// Reg is one register of the target description: the RSP register
// number (the "p N" address), the display name, and the width in bits.
// The order of Registers() is the 'g' block layout: the block offset of
// a register is the sum of the byte widths of its predecessors.
type Reg struct {
	name string
	num  int
	bits int
}

// NewReg - the register constructor: the struct is assembled only here.
func NewReg(name string, num, bits int) (Reg, error) {
	if name == "" {
		return Reg{}, fmt.Errorf("debug.NewReg: empty register name")
	}

	if num < 0 {
		return Reg{}, fmt.Errorf("debug.NewReg: register %q has a negative number %d", name, num)
	}

	if bits <= 0 || bits%8 != 0 {
		return Reg{}, fmt.Errorf(
			"debug.NewReg: register %q width %d is not a positive multiple of 8",
			name, bits,
		)
	}

	return Reg{
		name: name,
		num:  num,
		bits: bits,
	}, nil
}

// Name is the display name of the register.
func (r Reg) Name() string {
	return r.name
}

// Width is the register width in bytes (bits/8).
func (r Reg) Width() int {
	return r.bits / 8
}
