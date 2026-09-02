package arm64

// vOp - export surface of computed assembly operands for the syntax layer
// (asm/arm64): resolveOps lives there and builds these values via
// constructors; the fields stay unexported. The new* constructors and the
// legacy path consume vOp inside arch.

// VOp — computed assembly operand (see assemble.go).
type VOp = vOp

// VMem — computed memory operand.
type VMem = vMem

// ArmListReg — element of a register list ({ v0.16b, v1.16b }).
type ArmListReg = armListReg

// --- assembly (for the resolveOps of the syntax layer) ---

// VOpReg — a register (arr - the ".8b" suffix; lane/idx - the "v30[1]" suffix).
func VOpReg(reg, arr string, lane bool, idx int64) VOp {
	return VOp{kind: armOpReg, reg: reg, arr: arr, laneIdx: lane, num: idx}
}

// VOpImm — a number or name operand (condition/sysreg/prefetch hint; num
// and sym are mutually exclusive).
func VOpImm(num int64, sym string) VOp {
	return VOp{kind: armOpImm, num: num, sym: sym}
}

// VOpLit — a literal pool slot (num - the computed slot address).
func VOpLit(num int64) VOp {
	return VOp{kind: armOpLit, num: num}
}

// VOpFloat — a floating-point numeric literal (fmov imm).
func VOpFloat(f float64) VOp {
	return VOp{kind: armOpFloat, fval: f}
}

// VOpShift — a shift modifier (lsl/lsr/asr/ror #amt).
func VOpShift(shift string, amt int64, has bool) VOp {
	return VOp{kind: armOpShift, shift: shift, num: amt, hasAmt: has}
}

// VOpExtend — an extension modifier (uxtw/sxtw...[#amt]).
func VOpExtend(ext string, amt int64, has bool) VOp {
	return VOp{kind: armOpExtend, shift: ext, num: amt, hasAmt: has}
}

// VOpMem — a memory operand off(reg).
func VOpMem(m VMem) VOp {
	return VOp{kind: armOpMem, mem: &m}
}

// NewVMem — computed memory operand: [base, #off(!)] / [base, offReg
// opt #optAmt] / [base], #post.
func NewVMem(
	base string,
	off int64, hasOff bool,
	offReg string,
	opt string, optAmt int64, hasOpt bool,
	pre bool,
	post int64, hasPost bool,
) VMem {
	return VMem{
		base:    base,
		off:     off,
		hasOff:  hasOff,
		offReg:  offReg,
		opt:     opt,
		optAmt:  optAmt,
		hasOpt:  hasOpt,
		pre:     pre,
		post:    post,
		hasPost: hasPost,
	}
}

// VOpList — a register list ({ v0.16b, x0, x1 }).
func VOpList(regs []ArmListReg) VOp {
	return VOp{kind: armOpList, list: regs}
}

// NewArmListReg — a list element.
func NewArmListReg(reg, arr string) ArmListReg {
	return armListReg{reg: reg, arr: arr}
}

// --- reading (for the self-verify renderer of the syntax layer) ---

// Arr — the arrangement suffix ("" if absent).
func (o VOp) Arr() string {
	return o.arr
}

// LaneIdx — the register has a lane index suffix v30[1] (Num is its value).
func (o VOp) LaneIdx() bool {
	return o.laneIdx
}

// HasAmt — shift/extend has an amount set.
func (o VOp) HasAmt() bool {
	return o.hasAmt
}

// IsFloat — a floating-point numeric literal.
func (o VOp) IsFloat() bool {
	return o.kind == armOpFloat
}

// Float — the value of a floating-point literal.
func (o VOp) Float() float64 {
	return o.fval
}

// IsList — a register list.
func (o VOp) IsList() bool {
	return o.kind == armOpList
}

// List — elements of the register list.
func (o VOp) List() []ArmListReg {
	return o.list
}

// IsMem — a memory operand.
func (o VOp) IsMem() bool {
	return o.kind == armOpMem
}

// Mem — the memory operand (a copy).
func (o VOp) Mem() VMem {
	if o.mem == nil {
		return VMem{}
	}

	return *o.mem
}

// Base — the memory base register.
func (m VMem) Base() string {
	return m.base
}

// Off — the byte offset (HasOff tells whether it is set).
func (m VMem) Off() int64 {
	return m.off
}

// HasOff — a numeric offset is set.
func (m VMem) HasOff() bool {
	return m.hasOff
}

// OffReg — the index register ("" if absent).
func (m VMem) OffReg() string {
	return m.offReg
}

// Opt — the index register extension ("" if absent).
func (m VMem) Opt() string {
	return m.opt
}

// OptAmt — the extension amount (HasOpt tells whether it is set).
func (m VMem) OptAmt() int64 {
	return m.optAmt
}

// HasOpt — the extension amount is set.
func (m VMem) HasOpt() bool {
	return m.hasOpt
}

// Pre — the pre-index form [rn, #off]!.
func (m VMem) Pre() bool {
	return m.pre
}

// Post — the post-index imm ([rn], #post; HasPost tells whether it is set).
func (m VMem) Post() int64 {
	return m.post
}

// HasPost — post-index is set.
func (m VMem) HasPost() bool {
	return m.hasPost
}

// Reg — the register name of a list element.
func (r ArmListReg) Reg() string {
	return r.reg
}

// Arr — the arrangement suffix of a list element ("" if absent).
func (r ArmListReg) Arr() string {
	return r.arr
}

// ArrQSize — (Q, size) of an arrangement suffix (grammar validation).
func ArrQSize(arr string) (uint32, uint32, error) {
	return arrQSize(arr)
}
