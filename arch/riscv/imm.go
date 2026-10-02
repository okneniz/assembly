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

func (i Imm12) String() string {
	return fmt.Sprintf("%#x", i.v)
}

func (i Imm20) String() string {
	return fmt.Sprintf("%#x", i.v)
}

func (o Off) String() string {
	return fmt.Sprintf("%#x", o.v)
}
