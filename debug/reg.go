package debug

// Reg is one register of the target description: the RSP register
// number (the "p N" address), the display name, and the width in bits.
// The order of Registers() is the 'g' block layout: the block offset of
// a register is the sum of the byte widths of its predecessors.
type Reg struct {
	Name string
	Num  int
	Bits int
}

func NewReg(name string, num, bits int) Reg {
	return Reg{
		Name: name,
		Num:  num,
		Bits: bits,
	}
}

// Width is the register width in bytes (bits/8).
func (r Reg) Width() int {
	return r.Bits / 8
}
