package arm64

// The register operand: registers are encoded by number + class; the
// class decides how exactly register 31 is read (xzr/wzr/sp/wsp) - the
// choice stays with the user, there is no magic.

import (
	"fmt"
	"strconv"
)

// regClass - integer register class: width + meaning of register 31.
type regClass uint8

const (
	classX   regClass = iota // x0..x30
	classW                   // w0..w30
	classXZR                 // xzr (31st, zero)
	classWZR                 // wzr
	classSP                  // sp (31st, stack pointer)
	classWSP                 // wsp
)

// Reg — integer register. Register 31 is encoded via the class
// (XZR/WZR/SP/WSP): how exactly 31 is read in a particular instruction is
// decided by the context; the choice stays with the user - there is no magic.
type Reg struct {
	num   uint8 // 0..30; named 31st ones - 31
	class regClass
}

func newReg(num uint8, class regClass) Reg {
	return Reg{
		num:   num,
		class: class,
	}
}

// X — 64-bit register x0..x30.
func X(n int) (Reg, error) {
	if n < 0 || n > 30 {
		return Reg{}, fmt.Errorf(
			"arm64.X: register number %d is out of 0..30 (register 31 is XZR or SP)",
			n,
		)
	}

	return newReg(uint8(n), classX), nil
}

// W — 32-bit register w0..w30.
func W(n int) (Reg, error) {
	if n < 0 || n > 30 {
		return Reg{}, fmt.Errorf(
			"arm64.W: register number %d is out of 0..30 (register 31 is WZR or WSP)",
			n,
		)
	}

	return newReg(uint8(n), classW), nil
}

// Named 31st registers.
var (
	XZR = newReg(31, classXZR)
	WZR = newReg(31, classWZR)
	SP  = newReg(31, classSP)
	WSP = newReg(31, classWSP)
)

// Is64 — width of the class.
func (r Reg) Is64() bool {
	return r.class == classX || r.class == classXZR || r.class == classSP
}

func (r Reg) String() string {
	return r.name()
}

// bits - register number in the encoding (Rd/Rn/Rm field).
func (r Reg) bits() uint32 {
	return uint32(r.num)
}

// name - canonical name (string representation of internal instruction
// fields: "x5", "w3", "xzr", "sp", ...).
func (r Reg) name() string {
	switch r.class {
	case classX:
		return "x" + strconv.Itoa(int(r.num))
	case classW:
		return "w" + strconv.Itoa(int(r.num))
	case classXZR:
		return "xzr"
	case classWZR:
		return "wzr"
	case classSP:
		return "sp"
	default:
		return "wsp"
	}
}
