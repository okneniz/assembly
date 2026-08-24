package arm64

// Typed operands for instruction constructors: values are validated at
// creation - the constructor returns an error (panic is forbidden).
// Contextual constraints - what the 31st register means, offset alignment -
// are checked by instruction constructors via require* helpers (bottom of
// the file); everything else is the Encode error channel.

import (
	"fmt"
	"slices"
	"strconv"
)

// --- types -------------------------------------------------------------------

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

// Reg - integer register. Register 31 is encoded via the class
// (XZR/WZR/SP/WSP): how exactly 31 is read in a particular instruction is
// decided by the context; the choice stays with the user - there is no magic.
type Reg struct {
	num   uint8 // 0..30; named 31st ones - 31
	class regClass
}

func NewReg(num uint8, class regClass) Reg {
	return Reg{
		num:   num,
		class: class,
	}
}

// Imm12 — immediate 0..4095 (add/sub #imm).
type Imm12 struct {
	v uint32
}

// Imm16 — immediate 0..65535 (movz/movk, svc/brk).
type Imm16 struct {
	v uint32
}

// Imm6 - shift amount 0..63 (register operations).
type Imm6 struct {
	v uint32
}

// Hw - halfword position of movz/movk: encoded as lsl #hw*16.
type Hw uint8

const (
	Hw0 Hw = iota // no shift
	Hw1           // lsl #16
	Hw2           // lsl #32 (64-bit form only)
	Hw3           // lsl #48 (64-bit form only)
)

// Shift - kind of shift of the third operand of register operations.
type Shift uint8

const (
	LSL Shift = iota
	LSR
	ASR
	ROR
)

// Sh12 - shift of an add/sub immediate: none or lsl #12.
type Sh12 uint8

const (
	NoSh12 Sh12 = iota
	LSL12       // lsl #12
)

// Off - byte offset of a load/store (before scaling to imm12).
// Range and alignment depend on the access size - they are checked by the
// instruction constructor.
type Off int64

// --- constructors ------------------------------------------------------------

// X - 64-bit register x0..x30.
func X(n int) (Reg, error) {
	if n < 0 || n > 30 {
		return Reg{}, fmt.Errorf(
			"arm64.X: register number %d is out of 0..30 (register 31 is XZR or SP)",
			n,
		)
	}

	return NewReg(uint8(n), classX), nil
}

// W - 32-bit register w0..w30.
func W(n int) (Reg, error) {
	if n < 0 || n > 30 {
		return Reg{}, fmt.Errorf(
			"arm64.W: register number %d is out of 0..30 (register 31 is WZR or WSP)",
			n,
		)
	}

	return NewReg(uint8(n), classW), nil
}

// Named 31st registers.
var (
	XZR = NewReg(31, classXZR)
	WZR = NewReg(31, classWZR)
	SP  = NewReg(31, classSP)
	WSP = NewReg(31, classWSP)
)

// NewImm12 - validated value; error when out of range.
func NewImm12(v int64) (Imm12, error) {
	if v < 0 || v > 0xfff {
		return Imm12{}, fmt.Errorf("arm64.NewImm12: value %d is out of 0..4095", v)
	}

	return Imm12{uint32(v)}, nil
}

// NewImm16 - validated value; error when out of range.
func NewImm16(v int64) (Imm16, error) {
	if v < 0 || v > 0xffff {
		return Imm16{}, fmt.Errorf("arm64.NewImm16: value %d is out of 0..65535", v)
	}

	return Imm16{uint32(v)}, nil
}

// NewImm6 - validated value; error when out of range.
func NewImm6(v int64) (Imm6, error) {
	if v < 0 || v > 63 {
		return Imm6{}, fmt.Errorf("arm64.NewImm6: value %d is out of 0..63", v)
	}

	return Imm6{uint32(v)}, nil
}

// --- methods -------------------------------------------------------------------

// Is64 - width of the class.
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

// --- contextual checks (for instruction constructors) -------------------------

// requireClass - register of an allowed class, otherwise an error with a hint.
func requireClass(r Reg, instr, op, hint string, ok ...regClass) error {
	if slices.Contains(ok, r.class) {
		return nil
	}

	return fmt.Errorf("arm64.New%s: operand %s: %s does not fit - %s", instr, op, r.name(), hint)
}

// requireWidth - equal width of all registers of an instruction.
func requireWidth(instr string, regs ...Reg) error {
	for _, r := range regs[1:] {
		if r.Is64() != regs[0].Is64() {
			return fmt.Errorf("arm64.New%s: register widths must match: %s vs %s",
				instr, regs[0].name(), r.name())
		}
	}

	return nil
}

// requireHwW - the 32-bit form of movz/movk never has shift #32/#48.
func requireHwW(rd Reg, instr string, hw Hw) error {
	if !rd.Is64() && hw > Hw1 {
		return fmt.Errorf("arm64.New%s: the 32-bit form allows Hw0/Hw1, not %s", instr, hw)
	}

	return nil
}

// requireShift - constraints of the shifted form of add/sub: only
// lsl/lsr/asr (ror is a logical-family encoding, unallocated for add/sub);
// the 32-bit form limits the shift amount to 0..31 (imm6 >= 32 is
// unpredictable).
func requireShift(rd Reg, instr string, imm Imm6, sh Shift) error {
	if sh == ROR {
		return fmt.Errorf("arm64.New%s: ror is not allowed for add/sub (only lsl/lsr/asr)", instr)
	}

	if !rd.Is64() && imm.v > 31 {
		return fmt.Errorf(
			"arm64.New%s: the 32-bit form allows shift 0..31, not #%d",
			instr,
			imm.v,
		)
	}

	return nil
}

// requireOff - offset must be non-negative, aligned, and within imm12.
func requireOff(instr string, off Off, scale uint32) error {
	if off < 0 || uint64(off)&(uint64(1)<<scale-1) != 0 || uint64(off)>>scale > 0xfff {
		return fmt.Errorf(
			"arm64.New%s: offset %s is out of the form's range (0..0x%x, alignment %d)",
			instr,
			off,
			uint64(0xfff)<<scale,
			uint64(1)<<scale,
		)
	}

	return nil
}

// lsOperand - rt: x/w register (sp not allowed), rn: x/sp base.
func lsOperand(rt, rn Reg, instr string) error {
	if err := requireClass(rt, instr, "rt", "x/w register (register 31 in rt reads as zr)",
		classX, classW, classXZR, classWZR); err != nil {
		return err
	}

	return requireClass(rn, instr, "rn", "x register or SP (register 31 in the base reads as sp)",
		classX, classSP)
}
