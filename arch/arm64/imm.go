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

// Imm16 — immediate 0..65535 (movz/movk, svc/brk).
type Imm16 struct {
	v uint32
}

// Imm6 — shift amount 0..63 (register operations).
type Imm6 struct {
	v uint32
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

// Imm12 — validated value; error when out of range.
func (Builder) Imm12(v int64) (Imm12, error) {
	if v < 0 || v > 0xfff {
		return Imm12{}, fmt.Errorf("arm64.New().Imm12: value %d is out of 0..4095", v)
	}

	return Imm12{uint32(v)}, nil
}

// Imm16 — validated value; error when out of range.
func (Builder) Imm16(v int64) (Imm16, error) {
	if v < 0 || v > 0xffff {
		return Imm16{}, fmt.Errorf("arm64.New().Imm16: value %d is out of 0..65535", v)
	}

	return Imm16{uint32(v)}, nil
}

// Imm6 — validated value; error when out of range.
func (Builder) Imm6(v int64) (Imm6, error) {
	if v < 0 || v > 63 {
		return Imm6{}, fmt.Errorf("arm64.New().Imm6: value %d is out of 0..63", v)
	}

	return Imm6{uint32(v)}, nil
}

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
