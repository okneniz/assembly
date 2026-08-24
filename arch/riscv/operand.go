package riscv

// Typed operands for instruction constructors: values are validated
// at creation - the constructor returns an error (panic is forbidden).
// Registers are encoded by number; the canonical name is the
// ABI name (zero/ra/sp/t0/...), as the decoder prints it.

import "fmt"

// --- types -------------------------------------------------------------------

// Reg - an integer register x0..x31 (the name is an ABI alias of the number).
type Reg struct {
	num uint8
}

func NewReg(num uint8) Reg {
	return Reg{num: num}
}

// Imm12 - a signed I/S-type immediate: -2048..2047.
type Imm12 struct {
	v int64
}

// Imm20 - the lui U-type field: 0..0xfffff (the decoder reads it without
// sign extension - negative values do not wrap around).
type Imm20 struct {
	v int64
}

// Off - a load/store byte offset (I/S type, unscaled):
// -2048..2047.
type Off struct {
	v int64
}

// --- constructors ------------------------------------------------------------

// X - a register by number x0..x31.
func X(n int) (Reg, error) {
	if n < 0 || n > 31 {
		return Reg{}, fmt.Errorf("riscv.X: register number %d outside 0..31", n)
	}

	return NewReg(uint8(n)), nil
}

// Named ABI registers (others are available via X(n): t0 = X(5), a0 = X(10)...).
var (
	Zero = NewReg(0)
	Ra   = NewReg(1)
	Sp   = NewReg(2)
)

// NewImm12 - a validated value; an error when out of range.
func NewImm12(v int64) (Imm12, error) {
	if v < -2048 || v > 2047 {
		return Imm12{}, fmt.Errorf("riscv.NewImm12: value %d outside -2048..2047", v)
	}

	return Imm12{v}, nil
}

// NewImm20 - a validated value; an error when out of range.
func NewImm20(v int64) (Imm20, error) {
	if v < 0 || v > 0xfffff {
		return Imm20{}, fmt.Errorf("riscv.NewImm20: value %d outside 0..%d", v, 0xfffff)
	}

	return Imm20{v}, nil
}

// NewOff - a validated value; an error when out of range.
func NewOff(v int64) (Off, error) {
	if v < -2048 || v > 2047 {
		return Off{}, fmt.Errorf("riscv.NewOff: value %d outside -2048..2047", v)
	}

	return Off{v}, nil
}

// --- methods -------------------------------------------------------------------

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

func (i Imm12) String() string {
	return fmt.Sprintf("%#x", i.v)
}
func (i Imm20) String() string {
	return fmt.Sprintf("%#x", i.v)
}
func (o Off) String() string {
	return fmt.Sprintf("%#x", o.v)
}
