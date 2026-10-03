// Package arm64 - programs written directly in Go: a chain Program over
// the arch builders, labels resolved at assembly time. The counterpart of
// an .s source file: every chain method is a source line, a macro is an
// ordinary Go function returning *Program. Immediate operands are raw
// ints; validation errors are deferred to Assemble. The Go source
// position of each line (the debugger's address ↔ source map) comes from
// the position resolver injected with WithPos - there is no built-in
// caller detection.
// The chain is a producer over the unit output: every method deposits
// its record there (instructions as ready records, label-directed macros
// as deferred ones), and Assemble/AssembleLayout are the unit's resolve
// phase adapted to the chain's result shape.
package arm64

import (
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/unit"
)

// Program - a program being built: a chain over the unit output. Chain
// methods deposit one record each and return the program; nothing is
// encoded until Assemble. The stream mirror carries the chain's own
// guards (an instruction belongs to the text stream).
type Program struct {
	u       *unit.Unit
	errs    []error
	pos     func() unit.Pos
	b       arch.Builder
	stream  int  // the emitting stream: 0 text (default), 1 data
	sawData bool // any deposit ever went to the data stream
}

// New - a chain over the given output: the caller owns it (a compiler
// creates one output and passes it to every producer it composes - the
// chain does not know who else deposits there).
func New(u *unit.Unit) *Program {
	return &Program{u: u, pos: nopos}
}

// nopos - the default resolver: no position (a line reports one only
// after a resolver is injected with WithPos).
func nopos() unit.Pos {
	return unit.Pos{}
}

// Abs - abs.Arr vd, vn.
func (p *Program) Abs(rd, rn arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Abs(rd, rn, arr)
	return p.instrLine(pos, "abs", i, err)
}

// Adc - adc rd, rn, rm (add with carry).
func (p *Program) Adc(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Adc(rd, rn, rm)
	return p.instrLine(pos, "adc", i, err)
}

// Add - add.Arr vd, vn, vm.
func (p *Program) Add(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Add(rd, rn, rm, arr)
	return p.instrLine(pos, "add", i, err)
}

// AddExt - add rd, rn, rm[ext[#imm3]].
func (p *Program) AddExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.AddExt(rd, rn, rm, ext, imm3)
	return p.instrLine(pos, "add", i, err)
}

// AddImm - add rd, rn, #imm[, lsl #12].
func (p *Program) AddImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("add", err)
	}

	i, err := p.b.AddImm(rd, rn, v, sh)
	return p.instrLine(pos, "add", i, err)
}

// AddShift - add rd, rn, rm[, shift #imm].
func (p *Program) AddShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("add", err)
	}

	i, err := p.b.AddShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "add", i, err)
}

// Addp - addp.Arr vd, vn, vm.
func (p *Program) Addp(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Addp(rd, rn, rm, arr)
	return p.instrLine(pos, "addp", i, err)
}

// AddsExt - adds rd, rn, rm[ext[#imm3]].
func (p *Program) AddsExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.AddsExt(rd, rn, rm, ext, imm3)
	return p.instrLine(pos, "adds", i, err)
}

// AddsImm - adds rd, rn, #imm[, lsl #12] (the flag-setting add).
func (p *Program) AddsImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("adds", err)
	}

	i, err := p.b.AddsImm(rd, rn, v, sh)
	return p.instrLine(pos, "adds", i, err)
}

// AddsShift - adds rd, rn, rm[, shift #imm].
func (p *Program) AddsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("adds", err)
	}

	i, err := p.b.AddsShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "adds", i, err)
}

// laPair - adrp rd, page; add rd, :lo12:target.
func laPair(b arch.Builder, rd arch.Reg, target, pc int64) ([]unit.Resolved, error) {
	adrp, err := b.Adrp(rd, (target&^0xFFF-pc&^0xFFF)>>12)
	if err != nil {
		return nil, err
	}

	lo, err := b.Imm12(target & 0xFFF)
	if err != nil {
		return nil, err
	}

	add, err := b.AddImm(rd, rd, lo, arch.NoSh12)
	if err != nil {
		return nil, err
	}

	return []unit.Resolved{adrp, add}, nil
}

// Adr - load the address of a label into the register (pc-relative).
func (p *Program) Adr(rd arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "adr", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Adr(rd, int64(t)-int64(pc))
	})
}

// Adrp - load the 4KB page address of a label into the register
// (pc-relative; the offset is the page count between the instruction's
// page and the target's page).
func (p *Program) Adrp(rd arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "adrp", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Adrp(rd, (int64(t)&^0xFFF-int64(pc)&^0xFFF)>>12)
	})
}

// Aese - aese.16b vd, vn.
func (p *Program) Aese(rd, rn arch.VReg) *Program {
	pos := p.pos()
	i, err := p.b.Aese(rd, rn)
	return p.instrLine(pos, "aese", i, err)
}

// Aesmc - aesmc.16b vd, vn.
func (p *Program) Aesmc(rd, rn arch.VReg) *Program {
	pos := p.pos()
	i, err := p.b.Aesmc(rd, rn)
	return p.instrLine(pos, "aesmc", i, err)
}

// And - and.Arr vd, vn, vm.
func (p *Program) And(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.And(rd, rn, rm, arr)
	return p.instrLine(pos, "and", i, err)
}

// AndImm - and rd, rn, #imm (the bitmask immediate).
func (p *Program) AndImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.AndImm(rd, rn, imm)
	return p.instrLine(pos, "and", i, err)
}

// AndShift - and rd, rn, rm[, shift #imm].
func (p *Program) AndShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("and", err)
	}

	i, err := p.b.AndShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "and", i, err)
}

// AndsImm - ands rd, rn, #imm (the flag-setting and).
func (p *Program) AndsImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.AndsImm(rd, rn, imm)
	return p.instrLine(pos, "ands", i, err)
}

// AndsShift - ands rd, rn, rm[, shift #imm].
func (p *Program) AndsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("ands", err)
	}

	i, err := p.b.AndsShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "ands", i, err)
}

// Ascii - string data appended verbatim (no terminating zero), into the
// current stream.
func (p *Program) Ascii(s string) *Program {
	p.u.Ascii(p.pos(), s)
	return p
}

// AsrReg - asr rd, rn, rm (the register shift form).
func (p *Program) AsrReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.AsrReg(rd, rn, rm)
	return p.instrLine(pos, "asr", i, err)
}

// B - branch to a label.
func (p *Program) B(label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "b", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.B(int64(t) - int64(pc)), nil
	})
}

// Bcond - conditional branch to a label.
func (p *Program) Bcond(cond, label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "b."+cond, label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bcond(cond, int64(t)-int64(pc))
	})
}

// Bfm - bfm rd, rn, #immr, #imms.
func (p *Program) Bfm(rd, rn arch.Reg, immr, imms uint32) *Program {
	pos := p.pos()
	i, err := p.b.Bfm(rd, rn, immr, imms)
	return p.instrLine(pos, "bfm", i, err)
}

// Bic - bic.Arr vd, vn, vm.
func (p *Program) Bic(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Bic(rd, rn, rm, arr)
	return p.instrLine(pos, "bic", i, err)
}

// BicShift - bic rd, rn, rm[, shift #imm] (and-not).
func (p *Program) BicShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("bic", err)
	}

	i, err := p.b.BicShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "bic", i, err)
}

// BicsShift - bics rd, rn, rm[, shift #imm] (the flag-setting bic).
func (p *Program) BicsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("bics", err)
	}

	i, err := p.b.BicsShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "bics", i, err)
}

// Bif - bif.Arr vd, vn, vm.
func (p *Program) Bif(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Bif(rd, rn, rm, arr)
	return p.instrLine(pos, "bif", i, err)
}

// Bit - bit.Arr vd, vn, vm.
func (p *Program) Bit(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Bit(rd, rn, rm, arr)
	return p.instrLine(pos, "bit", i, err)
}

// Bl - branch with link to a label.
func (p *Program) Bl(label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "bl", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bl(int64(t) - int64(pc)), nil
	})
}

// Blr - blr rn (branch with link to register).
func (p *Program) Blr(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Blr(rn)
	return p.instrLine(pos, "blr", i, err)
}

// Br - br rn (branch to register).
func (p *Program) Br(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Br(rn)
	return p.instrLine(pos, "br", i, err)
}

// Brk - brk #imm.
func (p *Program) Brk(imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("brk", err)
	}

	return p.instrLine(pos, "brk", p.b.Brk(v), nil)
}

// Bsl - bsl.Arr vd, vn, vm.
func (p *Program) Bsl(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Bsl(rd, rn, rm, arr)
	return p.instrLine(pos, "bsl", i, err)
}

// Bss - a zero-fill reserve of the data stream: memory the kernel zeroes,
// no file bytes. Labels after it resolve past the reserved range.
func (p *Program) Bss(reserve int) *Program {
	if p.stream != 1 {
		return p.fail(".bss", errors.New(
			"a bss reserve belongs to the data stream",
		))
	}

	if reserve <= 0 {
		return p.fail(".bss", fmt.Errorf("reserve %d is not positive", reserve))
	}

	p.u.Bss(p.pos(), reserve)
	return p
}

// --- label-directed lines -------------------------------------------------------

// Build - materialize the program; deferred construction errors are
// returned alongside (an empty slice means the program is well-formed).
func (p *Program) Build() (*Binary, []error) {
	return &Binary{u: p.u, streams: p.sawData}, p.errs
}

// --- internals ---------------------------------------------------------------

// Bytes - raw data bytes, into the current stream.
func (p *Program) Bytes(b ...byte) *Program {
	p.u.Bytes(p.pos(), b...)
	return p
}

// Cbnz - branch to a label when the register is not zero.
func (p *Program) Cbnz(rt arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "cbnz", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Cbnz(rt, int64(t)-int64(pc))
	})
}

// Cbz - branch to a label when the register is zero.
func (p *Program) Cbz(rt arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "cbz", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Cbz(rt, int64(t)-int64(pc))
	})
}

// Ccmp - ccmp rn, rm, #nzcv, cond.
func (p *Program) Ccmp(rn, rm arch.Reg, nzcv uint32, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Ccmp(rn, rm, nzcv, cond)
	return p.instrLine(pos, "ccmp", i, err)
}

// Cls - cls rd, rn.
func (p *Program) Cls(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Cls(rd, rn)
	return p.instrLine(pos, "cls", i, err)
}

// Clz - clz rd, rn.
func (p *Program) Clz(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Clz(rd, rn)
	return p.instrLine(pos, "clz", i, err)
}

// Cmeq - cmeq.Arr vd, vn, vm.
func (p *Program) Cmeq(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Cmeq(rd, rn, rm, arr)
	return p.instrLine(pos, "cmeq", i, err)
}

// Cmge - cmge.Arr vd, vn, vm.
func (p *Program) Cmge(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Cmge(rd, rn, rm, arr)
	return p.instrLine(pos, "cmge", i, err)
}

// Cmtst - cmtst.Arr vd, vn, vm.
func (p *Program) Cmtst(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Cmtst(rd, rn, rm, arr)
	return p.instrLine(pos, "cmtst", i, err)
}

// --- SIMD ------------------------------------------------------------------------
// Cnt - cnt.Arr vd, vn.
func (p *Program) Cnt(rd, rn arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Cnt(rd, rn, arr)
	return p.instrLine(pos, "cnt", i, err)
}

// Csel - csel rd, rn, rm, cond.
func (p *Program) Csel(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csel(rd, rn, rm, cond)
	return p.instrLine(pos, "csel", i, err)
}

// Csinc - csinc rd, rn, rm, cond.
func (p *Program) Csinc(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csinc(rd, rn, rm, cond)
	return p.instrLine(pos, "csinc", i, err)
}

// Csinv - csinv rd, rn, rm, cond.
func (p *Program) Csinv(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csinv(rd, rn, rm, cond)
	return p.instrLine(pos, "csinv", i, err)
}

// Csneg - csneg rd, rn, rm, cond.
func (p *Program) Csneg(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csneg(rd, rn, rm, cond)
	return p.instrLine(pos, "csneg", i, err)
}

// --- shifts / bitfield ----------------------------------------------------------

// Data - emit into the data stream: the writable statics of the program.
// The stream split only takes effect in AssembleLayout (the flat
// Assemble rejects a program that has one).
func (p *Program) Data() *Program {
	p.stream = 1
	p.sawData = true
	p.u.Data()
	return p
}

// Dmb - dmb domain (the memory barrier of the domain).
func (p *Program) Dmb(domain arch.BarrierDomain) *Program {
	pos := p.pos()
	i, err := p.b.Dmb(domain)
	return p.instrLine(pos, "dmb", i, err)
}

// Dsb - dsb domain (the barrier completes before anything continues).
func (p *Program) Dsb(domain arch.BarrierDomain) *Program {
	pos := p.pos()
	i, err := p.b.Dsb(domain)
	return p.instrLine(pos, "dsb", i, err)
}

// Dup - dup.Arr vd, wn (DUP general).
func (p *Program) Dup(vd arch.VReg, wn arch.Reg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Dup(vd, wn, arr)
	return p.instrLine(pos, "dup", i, err)
}

// DupElem - dup.Arr vd, vn[idx] (DUP element).
func (p *Program) DupElem(rd, rn arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.DupElem(rd, rn, arr, idx)
	return p.instrLine(pos, "dup", i, err)
}

// DupScalar - mov <s|d>n, vn.s|d[idx] (the scalar DUP alias).
func (p *Program) DupScalar(rd, rn arch.VReg, elem string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.DupScalar(rd, rn, elem, idx)
	return p.instrLine(pos, "mov", i, err)
}

// Entry - the label the program starts at.
func (p *Program) Entry(name string) *Program {
	p.u.Entry(name)
	return p
}

// EonShift - eon rd, rn, rm[, shift #imm] (eor-not).
func (p *Program) EonShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("eon", err)
	}

	i, err := p.b.EonShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "eon", i, err)
}

// Eor - eor.Arr vd, vn, vm.
func (p *Program) Eor(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Eor(rd, rn, rm, arr)
	return p.instrLine(pos, "eor", i, err)
}

// EorImm - eor rd, rn, #imm (the bitmask immediate).
func (p *Program) EorImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.EorImm(rd, rn, imm)
	return p.instrLine(pos, "eor", i, err)
}

// EorShift - eor rd, rn, rm[, shift #imm].
func (p *Program) EorShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("eor", err)
	}

	i, err := p.b.EorShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "eor", i, err)
}

// Extr - extr rd, rn, rm, #lsb.
func (p *Program) Extr(rd, rn, rm arch.Reg, lsb int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(lsb)
	if err != nil {
		return p.fail("extr", err)
	}

	i, err := p.b.Extr(rd, rn, rm, v)
	return p.instrLine(pos, "extr", i, err)
}

// Fadd - fadd fd, fn, fm (double/single by the operand kind).
func (p *Program) Fadd(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fadd(rd, rn, rm)
	return p.instrLine(pos, "fadd", i, err)
}

// FcmlaElem - fcmla.Arr vd, vn, vm[idx], #rot.
func (p *Program) FcmlaElem(rd, rn, rm arch.VReg, arr string, idx, rot uint32) *Program {
	pos := p.pos()
	i, err := p.b.FcmlaElem(rd, rn, rm, arr, idx, rot)
	return p.instrLine(pos, "fcmla", i, err)
}

// --- floating point ------------------------------------------------------------

// Fcmp - fcmp fn, fm (sets the FP flags).
func (p *Program) Fcmp(rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcmp(rn, rm)
	return p.instrLine(pos, "fcmp", i, err)
}

// FcmpZero - fcmp fn, #0.0.
func (p *Program) FcmpZero(rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.FcmpZero(rn)
	return p.instrLine(pos, "fcmp", i, err)
}

// Fcvt - fcvt fd, fn (the s<->d width conversion).
func (p *Program) Fcvt(rd, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcvt(rd, rn)
	return p.instrLine(pos, "fcvt", i, err)
}

// Fcvtzs - fcvtzs wd|xd, fn (FP to signed integer).
func (p *Program) Fcvtzs(rd arch.Reg, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcvtzs(rd, rn)
	return p.instrLine(pos, "fcvtzs", i, err)
}

// Fcvtzu - fcvtzu wd|xd, fn (FP to unsigned integer).
func (p *Program) Fcvtzu(rd arch.Reg, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcvtzu(rd, rn)
	return p.instrLine(pos, "fcvtzu", i, err)
}

// Fdiv - fdiv fd, fn, fm.
func (p *Program) Fdiv(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fdiv(rd, rn, rm)
	return p.instrLine(pos, "fdiv", i, err)
}

// Fmadd - fmadd fd, fn, fm, fa (fd = fa + fn*fm).
func (p *Program) Fmadd(rd, rn, rm, ra arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmadd(rd, rn, rm, ra)
	return p.instrLine(pos, "fmadd", i, err)
}

// Fmax - fmax fd, fn, fm.
func (p *Program) Fmax(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmax(rd, rn, rm)
	return p.instrLine(pos, "fmax", i, err)
}

// Fmin - fmin fd, fn, fm.
func (p *Program) Fmin(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmin(rd, rn, rm)
	return p.instrLine(pos, "fmin", i, err)
}

// FmlaElem - fmla.Arr vd, vn, vm[idx] (by element).
func (p *Program) FmlaElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.FmlaElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "fmla", i, err)
}

// FmlsElem - fmls.Arr vd, vn, vm[idx] (by element).
func (p *Program) FmlsElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.FmlsElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "fmls", i, err)
}

// Fmov - fmov fd, fn (the register form between two FP registers).
func (p *Program) Fmov(rd, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmov(rd, rn)
	return p.instrLine(pos, "fmov", i, err)
}

// FmovFromGpr - fmov fd, xn | fmov sn, wn (bits from the integer file;
// fmov d0, xzr is the FP zero).
func (p *Program) FmovFromGpr(rd arch.FReg, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.FmovFromGpr(rd, rn)
	return p.instrLine(pos, "fmov", i, err)
}

// FmovImm - fmov fd, #imm (the VFP imm8 form).
func (p *Program) FmovImm(rd arch.FReg, imm float64) *Program {
	pos := p.pos()
	i, err := p.b.FmovImm(rd, imm)
	return p.instrLine(pos, "fmov", i, err)
}

// FmovToGpr - fmov xn, fd | fmov wn, sn (bits to the integer file).
func (p *Program) FmovToGpr(rd arch.Reg, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.FmovToGpr(rd, rn)
	return p.instrLine(pos, "fmov", i, err)
}

// Fmul - fmul fd, fn, fm.
func (p *Program) Fmul(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmul(rd, rn, rm)
	return p.instrLine(pos, "fmul", i, err)
}

// FmulElem - fmul.Arr vd, vn, vm[idx] (by element).
func (p *Program) FmulElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.FmulElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "fmul", i, err)
}

// FmulxElem - fmulx.Arr vd, vn, vm[idx] (by element).
func (p *Program) FmulxElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.FmulxElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "fmulx", i, err)
}

// Fneg - fneg fd, fn.
func (p *Program) Fneg(rd, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fneg(rd, rn)
	return p.instrLine(pos, "fneg", i, err)
}

// Fnmsub - fnmsub fd, fn, fm, fa (fd = -(fn*fm - fa)).
func (p *Program) Fnmsub(rd, rn, rm, ra arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fnmsub(rd, rn, rm, ra)
	return p.instrLine(pos, "fnmsub", i, err)
}

// Fsub - fsub fd, fn, fm.
func (p *Program) Fsub(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fsub(rd, rn, rm)
	return p.instrLine(pos, "fsub", i, err)
}

// Half - 16-bit little-endian values, into the current stream.
func (p *Program) Half(vs ...uint16) *Program {
	p.u.Half(p.pos(), vs...)
	return p
}

// Ins - mov.sz vd[idx], wn (INS general, the mov spelling).
func (p *Program) Ins(vd arch.VReg, idx uint32, wn arch.Reg, elem string) *Program {
	pos := p.pos()
	i, err := p.b.Ins(vd, idx, wn, elem)
	return p.instrLine(pos, "mov", i, err)
}

// InsElem - ins.sz vd[idx], vn[idx] (INS element).
func (p *Program) InsElem(rd, rn arch.VReg, elem string, idx, srcIdx uint32) *Program {
	pos := p.pos()
	i, err := p.b.InsElem(rd, rn, elem, idx, srcIdx)
	return p.instrLine(pos, "ins", i, err)
}

// Isb - isb (the instruction barrier; the full-system domain is its
// only form).
func (p *Program) Isb() *Program {
	pos := p.pos()
	i, err := p.b.Isb()
	return p.instrLine(pos, "isb", i, err)
}

// La - load the address of a label into the register: the adrp+add pair
// (a fixed 8 bytes; the page split is computed against the pair's own
// address, the low 12 bits come from the target - the arm64 twin of the
// riscv/loong64 La).
func (p *Program) La(rd arch.Reg, label string) *Program {
	if p.stream != 0 {
		return p.fail("la", errors.New("an instruction in the data stream"))
	}

	p.u.Sym(p.pos(), unit.NewPair("la", label, 8, func(t, pc uint64) ([]unit.Resolved, error) {
		return laPair(p.b, rd, int64(t), int64(pc))
	}))
	return p
}

// Label - define a label at the current position of the current stream.
func (p *Program) Label(name string) *Program {
	p.u.Label(name)
	return p
}

// Ldar - ldar rt, [rn].
func (p *Program) Ldar(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldar(rt, rn)
	return p.instrLine(pos, "ldar", i, err)
}

// Ldarb - ldarb rt, [rn].
func (p *Program) Ldarb(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldarb(rt, rn)
	return p.instrLine(pos, "ldarb", i, err)
}

// Ldaxr - ldaxr rt, [rn].
func (p *Program) Ldaxr(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldaxr(rt, rn)
	return p.instrLine(pos, "ldaxr", i, err)
}

// Ldaxrb - ldaxrb rt, [rn].
func (p *Program) Ldaxrb(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldaxrb(rt, rn)
	return p.instrLine(pos, "ldaxrb", i, err)
}

// Ldp - ldp rt, rt2, [rn, #off].
func (p *Program) Ldp(rt, rt2, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldp(rt, rt2, rn, arch.Off(off))
	return p.instrLine(pos, "ldp", i, err)
}

// Ldpsw - ldpsw rt, rt2, [rn, #off].
func (p *Program) Ldpsw(rt, rt2, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldpsw(rt, rt2, rn, arch.Off(off))
	return p.instrLine(pos, "ldpsw", i, err)
}

// --- stores ---------------------------------------------------------------------

// Ldr - ldr rt, [rn, #off].
func (p *Program) Ldr(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldr(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldr", i, err)
}

// LdrF - ldr st|dt, [rn, #off] (the FP/SIMD register form).
func (p *Program) LdrF(rt arch.FReg, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.LdrF(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldr", i, err)
}

// Ldrb - ldrb rt, [rn, #off].
func (p *Program) Ldrb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrb(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldrb", i, err)
}

// Ldrh - ldrh rt, [rn, #off].
func (p *Program) Ldrh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrh(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldrh", i, err)
}

// Ldrsb - ldrsb rt, [rn, #off].
func (p *Program) Ldrsb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrsb(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldrsb", i, err)
}

// Ldrsh - ldrsh rt, [rn, #off].
func (p *Program) Ldrsh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrsh(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldrsh", i, err)
}

// Ldrsw - ldrsw rt, [rn, #off].
func (p *Program) Ldrsw(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrsw(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldrsw", i, err)
}

// Ldur - ldur rt, [rn, #off] (unscaled).
func (p *Program) Ldur(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldur(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldur", i, err)
}

// Ldurb - ldurb rt, [rn, #off] (unscaled).
func (p *Program) Ldurb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldurb(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldurb", i, err)
}

// Ldurh - ldurh rt, [rn, #off] (unscaled).
func (p *Program) Ldurh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldurh(rt, rn, arch.Off(off))
	return p.instrLine(pos, "ldurh", i, err)
}

// LslReg - lsl rd, rn, rm (the register shift form).
func (p *Program) LslReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.LslReg(rd, rn, rm)
	return p.instrLine(pos, "lsl", i, err)
}

// LsrReg - lsr rd, rn, rm (the register shift form).
func (p *Program) LsrReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.LsrReg(rd, rn, rm)
	return p.instrLine(pos, "lsr", i, err)
}

// Madd - madd rd, rn, rm, ra (rd = ra + rn*rm).
func (p *Program) Madd(rd, rn, rm, ra arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Madd(rd, rn, rm, ra)
	return p.instrLine(pos, "madd", i, err)
}

// MlaElem - mla.Arr vd, vn, vm[idx] (by element).
func (p *Program) MlaElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.MlaElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "mla", i, err)
}

// MlsElem - mls.Arr vd, vn, vm[idx] (by element).
func (p *Program) MlsElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.MlsElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "mls", i, err)
}

// MovSimd - mov.8b|16b vd, vm (the ORR-vector mov alias).
func (p *Program) MovSimd(rd, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.MovSimd(rd, rm, arr)
	return p.instrLine(pos, "mov", i, err)
}

// Movk - movk rd, #imm, lsl #hw*16.
func (p *Program) Movk(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movk", err)
	}

	i, err := p.b.Movk(rd, v, hw)
	return p.instrLine(pos, "movk", i, err)
}

// Movn - movn rd, #imm[, lsl #hw*16] (movz's negative-immediate twin).
func (p *Program) Movn(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movn", err)
	}

	i, err := p.b.Movn(rd, v, hw)
	return p.instrLine(pos, "movn", i, err)
}

// Movz - movz rd, #imm[, lsl #hw*16].
func (p *Program) Movz(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movz", err)
	}

	i, err := p.b.Movz(rd, v, hw)
	return p.instrLine(pos, "movz", i, err)
}

// Mrs - mrs rd, sysreg.
func (p *Program) Mrs(rd arch.Reg, sysreg string) *Program {
	pos := p.pos()
	i, err := p.b.Mrs(rd, sysreg)
	return p.instrLine(pos, "mrs", i, err)
}

// Msr - msr sysreg, rt.
func (p *Program) Msr(sysreg string, rt arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Msr(sysreg, rt)
	return p.instrLine(pos, "msr", i, err)
}

// Msub - msub rd, rn, rm, ra (rd = ra - rn*rm).
func (p *Program) Msub(rd, rn, rm, ra arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Msub(rd, rn, rm, ra)
	return p.instrLine(pos, "msub", i, err)
}

// MulElem - mul.Arr vd, vn, vm[idx] (by element).
func (p *Program) MulElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.MulElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "mul", i, err)
}

// Nop - nop.
func (p *Program) Nop() *Program {
	pos := p.pos()
	return p.instrLine(pos, "nop", p.b.Nop(), nil)
}

// Not - not.Arr vd, vn.
func (p *Program) Not(rd, rn arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Not(rd, rn, arr)
	return p.instrLine(pos, "not", i, err)
}

// Orn - orn.Arr vd, vn, vm.
func (p *Program) Orn(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Orn(rd, rn, rm, arr)
	return p.instrLine(pos, "orn", i, err)
}

// OrnShift - orn rd, rn, rm[, shift #imm] (or-not).
func (p *Program) OrnShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("orn", err)
	}

	i, err := p.b.OrnShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "orn", i, err)
}

// Orr - orr.Arr vd, vn, vm.
func (p *Program) Orr(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Orr(rd, rn, rm, arr)
	return p.instrLine(pos, "orr", i, err)
}

// OrrImm - orr rd, rn, #imm (the bitmask immediate).
func (p *Program) OrrImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.OrrImm(rd, rn, imm)
	return p.instrLine(pos, "orr", i, err)
}

// OrrShift - orr rd, rn, rm[, shift #imm].
func (p *Program) OrrShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("orr", err)
	}

	i, err := p.b.OrrShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "orr", i, err)
}

// Prfm - prfm [rn].
func (p *Program) Prfm(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Prfm(rn)
	return p.instrLine(pos, "prfm", i, err)
}

// Quad - 64-bit little-endian values, into the current stream.
func (p *Program) Quad(vs ...uint64) *Program {
	p.u.Quad(p.pos(), vs...)
	return p
}

// Rbit - rbit rd, rn.
func (p *Program) Rbit(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rbit(rd, rn)
	return p.instrLine(pos, "rbit", i, err)
}

// RbitV - rbit.Arr vd, vn.
func (p *Program) RbitV(rd, rn arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.RbitV(rd, rn, arr)
	return p.instrLine(pos, "rbit", i, err)
}

// Ret - ret rn.
func (p *Program) Ret(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ret(rn)
	return p.instrLine(pos, "ret", i, err)
}

// --- loads ----------------------------------------------------------------------

// Rev - rev rd, rn.
func (p *Program) Rev(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rev(rd, rn)
	return p.instrLine(pos, "rev", i, err)
}

// Rev16 - rev16 rd, rn.
func (p *Program) Rev16(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rev16(rd, rn)
	return p.instrLine(pos, "rev16", i, err)
}

// Rev32 - rev32 rd, rn.
func (p *Program) Rev32(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rev32(rd, rn)
	return p.instrLine(pos, "rev32", i, err)
}

// --- control transfers ----------------------------------------------------------

// Rev32V - rev32.Arr vd, vn.
func (p *Program) Rev32V(rd, rn arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Rev32V(rd, rn, arr)
	return p.instrLine(pos, "rev32", i, err)
}

// RorReg - ror rd, rn, rm (the register shift form).
func (p *Program) RorReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.RorReg(rd, rn, rm)
	return p.instrLine(pos, "ror", i, err)
}

// Saddw - saddw.Arr vd, vn, vm.
func (p *Program) Saddw(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Saddw(rd, rn, rm, arr)
	return p.instrLine(pos, "saddw", i, err)
}

// Sbfm - sbfm rd, rn, #immr, #imms.
func (p *Program) Sbfm(rd, rn arch.Reg, immr, imms uint32) *Program {
	pos := p.pos()
	i, err := p.b.Sbfm(rd, rn, immr, imms)
	return p.instrLine(pos, "sbfm", i, err)
}

// Scvtf - scvtf fd, wn|xn (signed integer to FP).
func (p *Program) Scvtf(rd arch.FReg, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Scvtf(rd, rn)
	return p.instrLine(pos, "scvtf", i, err)
}

// Sdiv - sdiv rd, rn, rm.
func (p *Program) Sdiv(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Sdiv(rd, rn, rm)
	return p.instrLine(pos, "sdiv", i, err)
}

// Shl - shl.Arr vd, vn, #shift.
func (p *Program) Shl(rd, rn arch.VReg, arr string, shift uint32) *Program {
	pos := p.pos()
	i, err := p.b.Shl(rd, rn, arr, shift)
	return p.instrLine(pos, "shl", i, err)
}

// Smc - smc #imm (the secure-monitor call; PSCI rides it at #0).
func (p *Program) Smc(imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("smc", err)
	}

	return p.instrLine(pos, "smc", p.b.Smc(v), nil)
}

// SmlalElem - smlal{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) SmlalElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SmlalElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "smlal", i, err)
}

// SmlslElem - smlsl{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) SmlslElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SmlslElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "smlsl", i, err)
}

// Smov - smov wd, vn.sz[idx].
func (p *Program) Smov(wd arch.Reg, vn arch.VReg, elem string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.Smov(wd, vn, elem, idx)
	return p.instrLine(pos, "smov", i, err)
}

// Smulh - smulh rd, rn, rm (signed high 64 bits of the product).
func (p *Program) Smulh(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Smulh(rd, rn, rm)
	return p.instrLine(pos, "smulh", i, err)
}

// SmullElem - smull{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) SmullElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SmullElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "smull", i, err)
}

// SqdmlalElem - sqdmlal{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) SqdmlalElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqdmlalElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "sqdmlal", i, err)
}

// SqdmlslElem - sqdmlsl{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) SqdmlslElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqdmlslElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "sqdmlsl", i, err)
}

// SqdmulhElem - sqdmulh.Arr vd, vn, vm[idx] (by element).
func (p *Program) SqdmulhElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqdmulhElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "sqdmulh", i, err)
}

// SqdmullElem - sqdmull{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) SqdmullElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqdmullElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "sqdmull", i, err)
}

// SqrdmlahElem - sqrdmlah.Arr vd, vn, vm[idx] (by element).
func (p *Program) SqrdmlahElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqrdmlahElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "sqrdmlah", i, err)
}

// SqrdmlshElem - sqrdmlsh.Arr vd, vn, vm[idx] (by element).
func (p *Program) SqrdmlshElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqrdmlshElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "sqrdmlsh", i, err)
}

// SqrdmulhElem - sqrdmulh.Arr vd, vn, vm[idx] (by element).
func (p *Program) SqrdmulhElem(rd, rn, rm arch.VReg, arr string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.SqrdmulhElem(rd, rn, rm, arr, idx)
	return p.instrLine(pos, "sqrdmulh", i, err)
}

// Sqrshl - sqrshl.Arr vd, vn, vm.
func (p *Program) Sqrshl(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Sqrshl(rd, rn, rm, arr)
	return p.instrLine(pos, "sqrshl", i, err)
}

// Sri - sri.Arr vd, vn, #shift.
func (p *Program) Sri(rd, rn arch.VReg, arr string, shift uint32) *Program {
	pos := p.pos()
	i, err := p.b.Sri(rd, rn, arr, shift)
	return p.instrLine(pos, "sri", i, err)
}

// Sshr - sshr.Arr vd, vn, #shift.
func (p *Program) Sshr(rd, rn arch.VReg, arr string, shift uint32) *Program {
	pos := p.pos()
	i, err := p.b.Sshr(rd, rn, arr, shift)
	return p.instrLine(pos, "sshr", i, err)
}

// Ssubw - ssubw.Arr vd, vn, vm.
func (p *Program) Ssubw(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Ssubw(rd, rn, rm, arr)
	return p.instrLine(pos, "ssubw", i, err)
}

// Stlr - stlr rt, [rn].
func (p *Program) Stlr(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlr(rt, rn)
	return p.instrLine(pos, "stlr", i, err)
}

// Stlrb - stlrb rt, [rn].
func (p *Program) Stlrb(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlrb(rt, rn)
	return p.instrLine(pos, "stlrb", i, err)
}

// Stlxr - stlxr rs, rt, [rn].
func (p *Program) Stlxr(rs, rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlxr(rs, rt, rn)
	return p.instrLine(pos, "stlxr", i, err)
}

// Stlxrb - stlxrb rs, rt, [rn].
func (p *Program) Stlxrb(rs, rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlxrb(rs, rt, rn)
	return p.instrLine(pos, "stlxrb", i, err)
}

// Stp - stp rt, rt2, [rn, #off].
func (p *Program) Stp(rt, rt2, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Stp(rt, rt2, rn, arch.Off(off))
	return p.instrLine(pos, "stp", i, err)
}

// --- atomics --------------------------------------------------------------------

// Str - str rt, [rn, #off].
func (p *Program) Str(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Str(rt, rn, arch.Off(off))
	return p.instrLine(pos, "str", i, err)
}

// StrF - str st|dt, [rn, #off] (the FP/SIMD register form).
func (p *Program) StrF(rt arch.FReg, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.StrF(rt, rn, arch.Off(off))
	return p.instrLine(pos, "str", i, err)
}

// --- system ---------------------------------------------------------------------

// Strb - strb rt, [rn, #off].
func (p *Program) Strb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Strb(rt, rn, arch.Off(off))
	return p.instrLine(pos, "strb", i, err)
}

// Strh - strh rt, [rn, #off].
func (p *Program) Strh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Strh(rt, rn, arch.Off(off))
	return p.instrLine(pos, "strh", i, err)
}

// Stur - stur rt, [rn, #off] (unscaled).
func (p *Program) Stur(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Stur(rt, rn, arch.Off(off))
	return p.instrLine(pos, "stur", i, err)
}

// Sturb - sturb rt, [rn, #off] (unscaled).
func (p *Program) Sturb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Sturb(rt, rn, arch.Off(off))
	return p.instrLine(pos, "sturb", i, err)
}

// Sturh - sturh rt, [rn, #off] (unscaled).
func (p *Program) Sturh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Sturh(rt, rn, arch.Off(off))
	return p.instrLine(pos, "sturh", i, err)
}

// Stxrb - stxrb rs, rt, [rn].
func (p *Program) Stxrb(rs, rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stxrb(rs, rt, rn)
	return p.instrLine(pos, "stxrb", i, err)
}

// SubExt - sub rd, rn, rm[ext[#imm3]].
func (p *Program) SubExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.SubExt(rd, rn, rm, ext, imm3)
	return p.instrLine(pos, "sub", i, err)
}

// SubImm - sub rd, rn, #imm[, lsl #12].
func (p *Program) SubImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("sub", err)
	}

	i, err := p.b.SubImm(rd, rn, v, sh)
	return p.instrLine(pos, "sub", i, err)
}

// SubShift - sub rd, rn, rm[, shift #imm].
func (p *Program) SubShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("sub", err)
	}

	i, err := p.b.SubShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "sub", i, err)
}

// SubsExt - subs rd, rn, rm[ext[#imm3]].
func (p *Program) SubsExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.SubsExt(rd, rn, rm, ext, imm3)
	return p.instrLine(pos, "subs", i, err)
}

// SubsImm - subs rd, rn, #imm[, lsl #12] (the flag-setting sub).
func (p *Program) SubsImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("subs", err)
	}

	i, err := p.b.SubsImm(rd, rn, v, sh)
	return p.instrLine(pos, "subs", i, err)
}

// SubsShift - subs rd, rn, rm[, shift #imm].
func (p *Program) SubsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("subs", err)
	}

	i, err := p.b.SubsShift(rd, rn, rm, v, sh)
	return p.instrLine(pos, "subs", i, err)
}

// Svc - svc #imm (the canonical Darwin trap immediate is 0x80).
func (p *Program) Svc(imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("svc", err)
	}

	return p.instrLine(pos, "svc", p.b.Svc(v), nil)
}

// --- arithmetic -----------------------------------------------------------------

// Tbl - tbl.16b vd, { vn }, vm.
func (p *Program) Tbl(rd, rn, rm arch.VReg) *Program {
	pos := p.pos()
	i, err := p.b.Tbl(rd, rn, rm)
	return p.instrLine(pos, "tbl", i, err)
}

// Tbz - test a bit and branch to a label when it is zero (the register
// width is dictated by the bit number: 32..63 need an x register).
func (p *Program) Tbz(rt arch.Reg, bit uint32, label string) *Program {
	pos := p.pos()
	return p.branchLine(pos, "tbz", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Tbz(rt, bit, int64(t)-int64(pc))
	})
}

// --- instruction chain twins ----------------------------------------------------

// Text - emit into the text stream: the instructions and any read-only
// data placed before the first Data() call live here.
func (p *Program) Text() *Program {
	p.stream = 0
	p.u.Text()
	return p
}

// Uaddlv - uaddlv.Arr hN|sN|dN, vn.
func (p *Program) Uaddlv(rd, rn arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Uaddlv(rd, rn, arr)
	return p.instrLine(pos, "uaddlv", i, err)
}

// Uaddw - uaddw.Arr vd, vn, vm.
func (p *Program) Uaddw(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Uaddw(rd, rn, rm, arr)
	return p.instrLine(pos, "uaddw", i, err)
}

// Ubfm - ubfm rd, rn, #immr, #imms.
func (p *Program) Ubfm(rd, rn arch.Reg, immr, imms uint32) *Program {
	pos := p.pos()
	i, err := p.b.Ubfm(rd, rn, immr, imms)
	return p.instrLine(pos, "ubfm", i, err)
}

// Ucvtf - ucvtf fd, wn|xn (unsigned integer to FP).
func (p *Program) Ucvtf(rd arch.FReg, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ucvtf(rd, rn)
	return p.instrLine(pos, "ucvtf", i, err)
}

// Udiv - udiv rd, rn, rm.
func (p *Program) Udiv(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Udiv(rd, rn, rm)
	return p.instrLine(pos, "udiv", i, err)
}

// UmlalElem - umlal{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) UmlalElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.UmlalElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "umlal", i, err)
}

// UmlslElem - umlsl{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) UmlslElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.UmlslElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "umlsl", i, err)
}

// Umov - umov wd, vn.sz[idx].
func (p *Program) Umov(wd arch.Reg, vn arch.VReg, elem string, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.Umov(wd, vn, elem, idx)
	return p.instrLine(pos, "umov", i, err)
}

// Umulh - umulh rd, rn, rm (unsigned high 64 bits of the product).
func (p *Program) Umulh(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Umulh(rd, rn, rm)
	return p.instrLine(pos, "umulh", i, err)
}

// --- logical / conditional ------------------------------------------------------

// UmullElem - umull{,2}.Arr vd, vn, vm[idx] (by element, widening).
func (p *Program) UmullElem(rd, rn, rm arch.VReg, arr string, two bool, idx uint32) *Program {
	pos := p.pos()
	i, err := p.b.UmullElem(rd, rn, rm, arr, two, idx)
	return p.instrLine(pos, "umull", i, err)
}

// Ushr - ushr.Arr vd, vn, #shift.
func (p *Program) Ushr(rd, rn arch.VReg, arr string, shift uint32) *Program {
	pos := p.pos()
	i, err := p.b.Ushr(rd, rn, arr, shift)
	return p.instrLine(pos, "ushr", i, err)
}

// Usubw - usubw.Arr vd, vn, vm.
func (p *Program) Usubw(rd, rn, rm arch.VReg, arr string) *Program {
	pos := p.pos()
	i, err := p.b.Usubw(rd, rn, rm, arr)
	return p.instrLine(pos, "usubw", i, err)
}

// WithPos - the position resolver for all lines appended after this
// call: the file:line each chain call came from. There is no built-in
// caller detection - the resolver is always the user's; a caller-based
// one runs inside the chain method, where frame 0 is the resolver
// itself and frame 2 the code calling the chain (runtime.Caller(2)).
// A nil resolver is ignored.
func (p *Program) WithPos(pos func() unit.Pos) *Program {
	if pos != nil {
		p.pos = pos
	}

	return p
}

// Word - 32-bit little-endian values, into the current stream.
func (p *Program) Word(vs ...uint32) *Program {
	p.u.Word(p.pos(), vs...)
	return p
}

func (p *Program) branchLine(
	pos unit.Pos,
	src, label string,
	ctor func(target, pc uint64) (arch.Instr, error),
) *Program {
	if p.stream != 0 {
		return p.fail(src, errors.New("an instruction in the data stream"))
	}

	p.u.Sym(pos, unit.NewBranch(src, label, 4, func(t, pc uint64) (unit.Resolved, error) {
		return ctor(t, pc)
	}))
	return p
}

func (p *Program) fail(src string, err error) *Program {
	p.errs = append(p.errs, fmt.Errorf("%s: %w", src, err))
	return p
}

func (p *Program) instrLine(pos unit.Pos, src string, i arch.Instr, err error) *Program {
	if err != nil {
		return p.fail(src, err)
	}

	if p.stream != 0 {
		return p.fail(src, errors.New("an instruction in the data stream"))
	}

	p.u.Instr(pos, i, nil)
	return p
}
