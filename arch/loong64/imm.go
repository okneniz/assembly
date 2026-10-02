package loong64

// Immediate role types - the operand as the assembly writes it: values
// are validated at creation - the constructor returns an error (panic is
// forbidden).

import "fmt"

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
