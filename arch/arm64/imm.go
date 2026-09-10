package arm64

// Bounded scalar operand roles: immediate ranges are validated at
// creation - the constructor returns an error (panic is forbidden).
// Contextual constraints (offset alignment and range per access size)
// are checked by instruction constructors via the require* helpers.

import "fmt"

// Imm12 — immediate 0..4095 (add/sub #imm).
type Imm12 struct {
	v uint32
}

// newImm12 - the validating Imm12 constructor (the Builder method
// delegates here): the range check lives in the constructor.
func newImm12(v int64) (Imm12, error) {
	if v < 0 || v > 0xfff {
		return Imm12{}, fmt.Errorf("arm64.New().Imm12: value %d is out of 0..4095", v)
	}

	return Imm12{v: uint32(v)}, nil
}

// imm12Of - the word-bits form (the decoder): range-safe by
// construction.
func imm12Of(v uint32) Imm12 {
	return Imm12{v: v}
}

// Imm16 — immediate 0..65535 (movz/movk, svc/brk).
type Imm16 struct {
	v uint32
}

// newImm16 - the validating Imm16 constructor.
func newImm16(v int64) (Imm16, error) {
	if v < 0 || v > 0xffff {
		return Imm16{}, fmt.Errorf("arm64.New().Imm16: value %d is out of 0..65535", v)
	}

	return Imm16{v: uint32(v)}, nil
}

// imm16Of - the word-bits form (the decoder).
func imm16Of(v uint32) Imm16 {
	return Imm16{v: v}
}

// newImm6 - the validating Imm6 constructor.
func newImm6(v int64) (Imm6, error) {
	if v < 0 || v > 63 {
		return Imm6{}, fmt.Errorf("arm64.New().Imm6: value %d is out of 0..63", v)
	}

	return Imm6{v: uint32(v)}, nil
}

// imm6Of - the word-bits form (the decoder).
func imm6Of(v uint32) Imm6 {
	return Imm6{v: v}
}

// sh12Of - the word-bits form of Sh12 (the decoder).
func sh12Of(sh bool) Sh12 {
	if sh {
		return LSL12
	}

	return NoSh12
}

// shiftOf - the word-bits form of Shift (the decoder).
func shiftOf(kind uint32) Shift {
	return Shift(kind)
}

// Imm6 — shift amount 0..63 (register operations).
type Imm6 struct {
	v uint32
}

// shiftOfName - the Shift by its source name (the inverse of
// String; the string-operand layer).
func shiftOfName(name string) (Shift, error) {
	for i, n := range shiftNames {
		if n == name {
			return Shift(i), nil
		}
	}

	return LSL, fmt.Errorf("arm64: unknown shift %q", name)
}

// hwOf - the word-bits form of Hw (the decoder).
func hwOf(v uint32) Hw {
	return Hw(v)
}

// Hw — halfword position of movz/movk: encoded as lsl #hw*16.
type Hw uint8

const (
	Hw0 Hw = iota // no shift
	Hw1           // lsl #16
	Hw2           // lsl #32 (64-bit form only)
	Hw3           // lsl #48 (64-bit form only)
)

// Shift — kind of shift of the third operand of register operations.
type Shift uint8

const (
	LSL Shift = iota
	LSR
	ASR
	ROR
)

// Sh12 — shift of an add/sub immediate: none or lsl #12.
type Sh12 uint8

const (
	NoSh12 Sh12 = iota
	LSL12       // lsl #12
)

// Off — byte offset of a load/store (before scaling to imm12).
// Range and alignment depend on the access size - they are checked by the
// instruction constructor.
type Off int64

func (i Imm12) String() string {
	return fmt.Sprintf("#0x%x", i.v)
}

func (i Imm16) String() string {
	return fmt.Sprintf("#0x%x", i.v)
}

func (i Imm6) String() string {
	return fmt.Sprintf("#%d", i.v)
}

func (h Hw) String() string {
	return fmt.Sprintf("lsl #%d", uint32(h)*16)
}

func (s Shift) String() string {
	return shiftNames[s]
}

func (s Sh12) String() string {
	if s == LSL12 {
		return "lsl #12"
	}

	return ""
}

func (o Off) String() string {
	return fmt.Sprintf("#%#x", int64(o))
}

// Imm12 — validated value; error when out of range.
func (Builder) Imm12(v int64) (Imm12, error) {
	return newImm12(v)
}

// Imm16 — validated value; error when out of range.
func (Builder) Imm16(v int64) (Imm16, error) {
	return newImm16(v)
}

// Imm6 — validated value; error when out of range.
func (Builder) Imm6(v int64) (Imm6, error) {
	return newImm6(v)
}
