package loong64

// Typed operands for instruction constructors: values are validated
// at creation - the constructor returns an error (panic is forbidden).
// Instructions store register NUMBERS (a byte each, not names); the
// canonical $-name is only produced at print time.

import "fmt"

// Reg - an integer register $r0..$r31 (the name is an ABI alias of the number).
type Reg struct {
	num uint8
}

func NewReg(num uint8) Reg {
	return Reg{
		num: num,
	}
}

// R - a register by number $r0..$r31.
func R(n int) (Reg, error) {
	if n < 0 || n > 31 {
		return Reg{}, fmt.Errorf("loong64.R: register number %d outside 0..31", n)
	}

	return NewReg(uint8(n)), nil
}

// Named ABI registers (others are available via R(n): t0 = R(12), a0 = R(4)...).
var (
	Zero = NewReg(0)
	Ra   = NewReg(1)
	Tp   = NewReg(2)
	Sp   = NewReg(3)
	Fp   = NewReg(22)
)

// --- immediate role types (the operand as the assembly writes it) -------------

// Imm12 - a signed si12 immediate: -2048..2047 (ALU immediates, ld/st byte
// offsets).
type Imm12 struct {
	v int64
}

// UImm12 - an unsigned ui12 immediate: 0..4095 (andi/ori/xori).
type UImm12 struct {
	v int64
}

// Imm14 - a word-scaled si14 byte offset (ldptr/stptr): -16380..16380, a
// multiple of 4.
type Imm14 struct {
	v int64
}

// Imm16 - a signed si16 immediate: -32768..32767 (addu16i.d).
type Imm16 struct {
	v int64
}

// Off16 - a word-scaled si16 byte offset (branches, jirl): -131068..131068,
// a multiple of 4.
type Off16 struct {
	v int64
}

// Imm20 - a signed si20 immediate: -524288..524287 (lu12i.w and the
// pcaddi family).
type Imm20 struct {
	v int64
}

// UImm5 - an unsigned 5-bit immediate: 0..31 (.w shift amounts, bit field
// bounds).
type UImm5 struct {
	v int64
}

// UImm2 - an unsigned 2-bit immediate: 0..3 (the bytepick.w index).
type UImm2 struct {
	v int64
}

// UImm3 - an unsigned 3-bit immediate: 0..7 (the bytepick.d index).
type UImm3 struct {
	v int64
}

// UImm6 - an unsigned 6-bit immediate: 0..63 (.d shift amounts).
type UImm6 struct {
	v int64
}

// Shift3 - an alsl shift amount: 1..4 (encoded shifted by one).
type Shift3 struct {
	v int64
}

// UImm8 - an unsigned 8-bit immediate: 0..255 (lddir/ldpte).
type UImm8 struct {
	v int64
}

// UImm14 - an unsigned 14-bit CSR number: 0..16383.
type UImm14 struct {
	v int64
}

// Code15 - an unsigned 15-bit code: 0..32767 (break/syscall/dbar/ibar).
type Code15 struct {
	v int64
}

func newImm(v, lo, hi int64, what string) (int64, error) {
	if v < lo || v > hi {
		return 0, fmt.Errorf("loong64: %s %d outside %d..%d", what, v, lo, hi)
	}

	return v, nil
}

// NewImm12 - a validated value; an error when out of range.
func NewImm12(v int64) (Imm12, error) {
	x, err := newImm(v, -2048, 2047, "si12 value")
	if err != nil {
		return Imm12{}, err
	}

	return Imm12{v: x}, nil
}

// NewUImm12 - a validated value; an error when out of range.
func NewUImm12(v int64) (UImm12, error) {
	x, err := newImm(v, 0, 4095, "ui12 value")
	if err != nil {
		return UImm12{}, err
	}

	return UImm12{v: x}, nil
}

// NewImm14 - a validated word-aligned value; an error otherwise.
func NewImm14(v int64) (Imm14, error) {
	if v%4 != 0 {
		return Imm14{}, fmt.Errorf("loong64: si14 byte offset %d not word-aligned", v)
	}

	x, err := newImm(v, -16380, 16380, "si14 byte offset")
	if err != nil {
		return Imm14{}, err
	}

	return Imm14{v: x}, nil
}

// NewImm16 - a validated value; an error when out of range.
func NewImm16(v int64) (Imm16, error) {
	x, err := newImm(v, -32768, 32767, "si16 value")
	if err != nil {
		return Imm16{}, err
	}

	return Imm16{v: x}, nil
}

// NewOff16 - a validated word-aligned value; an error otherwise.
func NewOff16(v int64) (Off16, error) {
	if v%4 != 0 {
		return Off16{}, fmt.Errorf("loong64: si16 byte offset %d not word-aligned", v)
	}

	x, err := newImm(v, -131068, 131068, "si16 byte offset")
	if err != nil {
		return Off16{}, err
	}

	return Off16{v: x}, nil
}

// NewImm20 - a validated value; an error when out of range.
func NewImm20(v int64) (Imm20, error) {
	x, err := newImm(v, -524288, 524287, "si20 value")
	if err != nil {
		return Imm20{}, err
	}

	return Imm20{v: x}, nil
}

// NewUImm5 - a validated value; an error when out of range.
func NewUImm5(v int64) (UImm5, error) {
	x, err := newImm(v, 0, 31, "ui5 value")
	if err != nil {
		return UImm5{}, err
	}

	return UImm5{v: x}, nil
}

// NewUImm2 - a validated value; an error when out of range.
func NewUImm2(v int64) (UImm2, error) {
	x, err := newImm(v, 0, 3, "ui2 value")
	if err != nil {
		return UImm2{}, err
	}

	return UImm2{v: x}, nil
}

// NewUImm3 - a validated value; an error when out of range.
func NewUImm3(v int64) (UImm3, error) {
	x, err := newImm(v, 0, 7, "ui3 value")
	if err != nil {
		return UImm3{}, err
	}

	return UImm3{v: x}, nil
}

// NewUImm6 - a validated value; an error when out of range.
func NewUImm6(v int64) (UImm6, error) {
	x, err := newImm(v, 0, 63, "ui6 value")
	if err != nil {
		return UImm6{}, err
	}

	return UImm6{v: x}, nil
}

// NewShift3 - a validated alsl shift amount (1..4).
func NewShift3(v int64) (Shift3, error) {
	x, err := newImm(v, 1, 4, "alsl shift amount")
	if err != nil {
		return Shift3{}, err
	}

	return Shift3{v: x}, nil
}

// NewUImm8 - a validated value; an error when out of range.
func NewUImm8(v int64) (UImm8, error) {
	x, err := newImm(v, 0, 255, "ui8 value")
	if err != nil {
		return UImm8{}, err
	}

	return UImm8{v: x}, nil
}

// NewUImm14 - a validated CSR number; an error when out of range.
func NewUImm14(v int64) (UImm14, error) {
	x, err := newImm(v, 0, 16383, "csr number")
	if err != nil {
		return UImm14{}, err
	}

	return UImm14{v: x}, nil
}

// NewCode15 - a validated code; an error when out of range.
func NewCode15(v int64) (Code15, error) {
	x, err := newImm(v, 0, 32767, "code")
	if err != nil {
		return Code15{}, err
	}

	return Code15{v: x}, nil
}

// Num - the register number (rd/rj/rk field: 0..31).
func (r Reg) Num() uint8 {
	return r.num
}

// String - the printed form ("$zero", "$a0", "$r21", ...).
func (r Reg) String() string {
	return laRegName(r.num)
}

// Val - the immediate value.
func (i Imm12) Val() int64 {
	return i.v
}

func (i UImm12) Val() int64 {
	return i.v
}

func (i Imm14) Val() int64 {
	return i.v
}

func (i Imm16) Val() int64 {
	return i.v
}

func (i Off16) Val() int64 {
	return i.v
}

func (i Imm20) Val() int64 {
	return i.v
}

func (i UImm5) Val() int64 {
	return i.v
}

func (i UImm2) Val() int64 {
	return i.v
}

func (i UImm3) Val() int64 {
	return i.v
}

func (i UImm6) Val() int64 {
	return i.v
}

func (i Shift3) Val() int64 {
	return i.v
}

func (i UImm8) Val() int64 {
	return i.v
}

func (i UImm14) Val() int64 {
	return i.v
}

func (i Code15) Val() int64 {
	return i.v
}
