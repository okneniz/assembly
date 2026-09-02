package riscv

// Bounded scalar operand roles for instruction constructors: values are
// validated at creation - the constructor returns an error (panic is
// forbidden).

import "fmt"

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

// Imm12 - a validated value; an error when out of range.
func (Builder) Imm12(v int64) (Imm12, error) {
	if v < -2048 || v > 2047 {
		return Imm12{}, fmt.Errorf("riscv.New().Imm12: value %d outside -2048..2047", v)
	}

	return Imm12{v}, nil
}

// Imm20 - a validated value; an error when out of range.
func (Builder) Imm20(v int64) (Imm20, error) {
	if v < 0 || v > 0xfffff {
		return Imm20{}, fmt.Errorf("riscv.New().Imm20: value %d outside 0..%d", v, 0xfffff)
	}

	return Imm20{v}, nil
}

// Off - a validated value; an error when out of range.
func (Builder) Off(v int64) (Off, error) {
	if v < -2048 || v > 2047 {
		return Off{}, fmt.Errorf("riscv.New().Off: value %d outside -2048..2047", v)
	}

	return Off{v}, nil
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
