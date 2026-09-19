// Package arm64 - programs written directly in Go: a chain Program over
// the arch builders, labels resolved at assembly time. The counterpart of
// an .s source file: every chain method is a source line, a macro is an
// ordinary Go function returning *Program. Immediate operands are raw
// ints; validation errors are deferred to Assemble. The Go source
// position of each line (the debugger's address ↔ source map) comes from
// the position resolver injected with WithPos - there is no built-in
// caller detection.
package arm64

import (
	"encoding/binary"
	"errors"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

// Program - a program being built: a sequence of lines (instructions,
// label-directed branches, labels, data). Chain methods append one line
// each and return the program; nothing is encoded until Assemble.
type Program struct {
	lines  []line
	entry  string
	errs   []error
	pos    func() prog.Pos
	b      arch.Builder
	stream int // the emitting stream: 0 text (default), 1 data
}

// New - an empty program.
func New() *Program {
	return &Program{pos: nopos}
}

// nopos - the default resolver: no position (a line reports one only
// after a resolver is injected with WithPos).
func nopos() prog.Pos {
	return prog.Pos{}
}

// WithPos - the position resolver for all lines appended after this
// call: the file:line each chain call came from. There is no built-in
// caller detection - the resolver is always the user's; a caller-based
// one runs inside the chain method, where frame 0 is the resolver
// itself and frame 2 the code calling the chain (runtime.Caller(2)).
// A nil resolver is ignored.
func (p *Program) WithPos(pos func() prog.Pos) *Program {
	if pos != nil {
		p.pos = pos
	}

	return p
}

// Label - define a label at the current position of the current stream.
func (p *Program) Label(name string) *Program {
	p.lines = append(p.lines, newLabelLine(name, p.stream, p.pos()))
	return p
}

// Entry - the label the program starts at.
func (p *Program) Entry(name string) *Program {
	p.entry = name
	return p
}

// Text - emit into the text stream: the instructions and any read-only
// data placed before the first Data() call live here.
func (p *Program) Text() *Program {
	p.stream = 0
	return p
}

// Data - emit into the data stream: the writable statics of the program.
// The stream split only takes effect in AssembleLayout (the flat
// Assemble rejects a program that has one).
func (p *Program) Data() *Program {
	p.stream = 1
	return p
}

// Ascii - string data appended verbatim (no terminating zero), into the
// current stream.
func (p *Program) Ascii(s string) *Program {
	p.lines = append(p.lines, newDataLine([]byte(s), ".ascii", p.stream, p.pos()))
	return p
}

// Bytes - raw data bytes, into the current stream.
func (p *Program) Bytes(b ...byte) *Program {
	p.lines = append(p.lines, newDataLine(b, ".byte", p.stream, p.pos()))
	return p
}

// Half - 16-bit little-endian values, into the current stream.
func (p *Program) Half(vs ...uint16) *Program {
	b := make([]byte, 2*len(vs))
	for i, v := range vs {
		binary.LittleEndian.PutUint16(b[2*i:], v)
	}

	p.lines = append(p.lines, newDataLine(b, ".half", p.stream, p.pos()))
	return p
}

// Word - 32-bit little-endian values, into the current stream.
func (p *Program) Word(vs ...uint32) *Program {
	b := make([]byte, 4*len(vs))
	for i, v := range vs {
		binary.LittleEndian.PutUint32(b[4*i:], v)
	}

	p.lines = append(p.lines, newDataLine(b, ".word", p.stream, p.pos()))
	return p
}

// Quad - 64-bit little-endian values, into the current stream.
func (p *Program) Quad(vs ...uint64) *Program {
	b := make([]byte, 8*len(vs))
	for i, v := range vs {
		binary.LittleEndian.PutUint64(b[8*i:], v)
	}

	p.lines = append(p.lines, newDataLine(b, ".quad", p.stream, p.pos()))
	return p
}

// Bss - a zero-fill reserve of the data stream: memory the kernel zeroes,
// no file bytes. Labels after it resolve past the reserved range.
func (p *Program) Bss(reserve int) *Program {
	pos := p.pos()
	if p.stream != 1 {
		return p.fail(".bss", errors.New(
			"a bss reserve belongs to the data stream",
		))
	}

	if reserve <= 0 {
		return p.fail(".bss", fmt.Errorf("reserve %d is not positive", reserve))
	}

	p.lines = append(p.lines, newBssLine(reserve, pos))
	return p
}

// --- label-directed lines -------------------------------------------------------

// B - branch to a label.
func (p *Program) B(label string) *Program {
	pos := p.pos()
	return p.branchLine("b", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.B(int64(t) - int64(pc)), nil
	}, pos)
}

// Bl - branch with link to a label.
func (p *Program) Bl(label string) *Program {
	pos := p.pos()
	return p.branchLine("bl", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bl(int64(t) - int64(pc)), nil
	}, pos)
}

// Bcond - conditional branch to a label.
func (p *Program) Bcond(cond, label string) *Program {
	pos := p.pos()
	return p.branchLine("b."+cond, label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bcond(cond, int64(t)-int64(pc))
	}, pos)
}

// Cbz - branch to a label when the register is zero.
func (p *Program) Cbz(rt arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("cbz", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Cbz(rt, int64(t)-int64(pc))
	}, pos)
}

// Cbnz - branch to a label when the register is not zero.
func (p *Program) Cbnz(rt arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("cbnz", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Cbnz(rt, int64(t)-int64(pc))
	}, pos)
}

// Adr - load the address of a label into the register (pc-relative).
func (p *Program) Adr(rd arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("adr", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Adr(rd, int64(t)-int64(pc))
	}, pos)
}

// Adrp - load the 4KB page address of a label into the register
// (pc-relative; the offset is the page count between the instruction's
// page and the target's page).
func (p *Program) Adrp(rd arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("adrp", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Adrp(rd, (int64(t)&^0xFFF-int64(pc)&^0xFFF)>>12)
	}, pos)
}

// La - load the address of a label into the register: the adrp+add pair
// (a fixed 8 bytes; the page split is computed against the pair's own
// address, the low 12 bits come from the target - the arm64 twin of the
// riscv/loong64 La).
func (p *Program) La(rd arch.Reg, label string) *Program {
	pos := p.pos()
	p.lines = append(p.lines, newLaLine(label, func(t, pc uint64) ([]arch.Instr, error) {
		return laPair(p.b, rd, int64(t), int64(pc))
	}, pos))
	return p
}

// laPair - adrp rd, page; add rd, rd, :lo12:target.
func laPair(b arch.Builder, rd arch.Reg, target, pc int64) ([]arch.Instr, error) {
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

	return []arch.Instr{adrp, add}, nil
}

// Tbz - test a bit and branch to a label when it is zero (the register
// width is dictated by the bit number: 32..63 need an x register).
func (p *Program) Tbz(rt arch.Reg, bit uint32, label string) *Program {
	pos := p.pos()
	return p.branchLine("tbz", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Tbz(rt, bit, int64(t)-int64(pc))
	}, pos)
}

// --- instruction chain twins ----------------------------------------------------

// Movz - movz rd, #imm[, lsl #hw*16].
func (p *Program) Movz(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movz", err)
	}

	i, err := p.b.Movz(rd, v, hw)
	return p.instrLine("movz", i, err, pos)
}

// Movk - movk rd, #imm, lsl #hw*16.
func (p *Program) Movk(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movk", err)
	}

	i, err := p.b.Movk(rd, v, hw)
	return p.instrLine("movk", i, err, pos)
}

// Movn - movn rd, #imm[, lsl #hw*16] (movz's negative-immediate twin).
func (p *Program) Movn(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movn", err)
	}

	i, err := p.b.Movn(rd, v, hw)
	return p.instrLine("movn", i, err, pos)
}

// Svc - svc #imm (the canonical Darwin trap immediate is 0x80).
func (p *Program) Svc(imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("svc", err)
	}

	return p.instrLine("svc", p.b.Svc(v), nil, pos)
}

// --- arithmetic -----------------------------------------------------------------

// Adc - adc rd, rn, rm (add with carry).
func (p *Program) Adc(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Adc(rd, rn, rm)
	return p.instrLine("adc", i, err, pos)
}

// AddImm - add rd, rn, #imm[, lsl #12].
func (p *Program) AddImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("add", err)
	}

	i, err := p.b.AddImm(rd, rn, v, sh)
	return p.instrLine("add", i, err, pos)
}

// AddShift - add rd, rn, rm[, shift #imm].
func (p *Program) AddShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("add", err)
	}

	i, err := p.b.AddShift(rd, rn, rm, v, sh)
	return p.instrLine("add", i, err, pos)
}

// AddExt - add rd, rn, rm[ext[#imm3]].
func (p *Program) AddExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.AddExt(rd, rn, rm, ext, imm3)
	return p.instrLine("add", i, err, pos)
}

// AddsImm - adds rd, rn, #imm[, lsl #12] (the flag-setting add).
func (p *Program) AddsImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("adds", err)
	}

	i, err := p.b.AddsImm(rd, rn, v, sh)
	return p.instrLine("adds", i, err, pos)
}

// AddsShift - adds rd, rn, rm[, shift #imm].
func (p *Program) AddsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("adds", err)
	}

	i, err := p.b.AddsShift(rd, rn, rm, v, sh)
	return p.instrLine("adds", i, err, pos)
}

// AddsExt - adds rd, rn, rm[ext[#imm3]].
func (p *Program) AddsExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.AddsExt(rd, rn, rm, ext, imm3)
	return p.instrLine("adds", i, err, pos)
}

// SubImm - sub rd, rn, #imm[, lsl #12].
func (p *Program) SubImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("sub", err)
	}

	i, err := p.b.SubImm(rd, rn, v, sh)
	return p.instrLine("sub", i, err, pos)
}

// SubShift - sub rd, rn, rm[, shift #imm].
func (p *Program) SubShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("sub", err)
	}

	i, err := p.b.SubShift(rd, rn, rm, v, sh)
	return p.instrLine("sub", i, err, pos)
}

// SubExt - sub rd, rn, rm[ext[#imm3]].
func (p *Program) SubExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.SubExt(rd, rn, rm, ext, imm3)
	return p.instrLine("sub", i, err, pos)
}

// SubsImm - subs rd, rn, #imm[, lsl #12] (the flag-setting sub).
func (p *Program) SubsImm(rd, rn arch.Reg, imm int64, sh arch.Sh12) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("subs", err)
	}

	i, err := p.b.SubsImm(rd, rn, v, sh)
	return p.instrLine("subs", i, err, pos)
}

// SubsShift - subs rd, rn, rm[, shift #imm].
func (p *Program) SubsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("subs", err)
	}

	i, err := p.b.SubsShift(rd, rn, rm, v, sh)
	return p.instrLine("subs", i, err, pos)
}

// SubsExt - subs rd, rn, rm[ext[#imm3]].
func (p *Program) SubsExt(rd, rn, rm arch.Reg, ext string, imm3 uint32) *Program {
	pos := p.pos()
	i, err := p.b.SubsExt(rd, rn, rm, ext, imm3)
	return p.instrLine("subs", i, err, pos)
}

// Madd - madd rd, rn, rm, ra (rd = ra + rn*rm).
func (p *Program) Madd(rd, rn, rm, ra arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Madd(rd, rn, rm, ra)
	return p.instrLine("madd", i, err, pos)
}

// Msub - msub rd, rn, rm, ra (rd = ra - rn*rm).
func (p *Program) Msub(rd, rn, rm, ra arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Msub(rd, rn, rm, ra)
	return p.instrLine("msub", i, err, pos)
}

// Sdiv - sdiv rd, rn, rm.
func (p *Program) Sdiv(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Sdiv(rd, rn, rm)
	return p.instrLine("sdiv", i, err, pos)
}

// Udiv - udiv rd, rn, rm.
func (p *Program) Udiv(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Udiv(rd, rn, rm)
	return p.instrLine("udiv", i, err, pos)
}

// Smulh - smulh rd, rn, rm (signed high 64 bits of the product).
func (p *Program) Smulh(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Smulh(rd, rn, rm)
	return p.instrLine("smulh", i, err, pos)
}

// Umulh - umulh rd, rn, rm (unsigned high 64 bits of the product).
func (p *Program) Umulh(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Umulh(rd, rn, rm)
	return p.instrLine("umulh", i, err, pos)
}

// --- logical / conditional ------------------------------------------------------

// AndImm - and rd, rn, #imm (the bitmask immediate).
func (p *Program) AndImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.AndImm(rd, rn, imm)
	return p.instrLine("and", i, err, pos)
}

// AndShift - and rd, rn, rm[, shift #imm].
func (p *Program) AndShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("and", err)
	}

	i, err := p.b.AndShift(rd, rn, rm, v, sh)
	return p.instrLine("and", i, err, pos)
}

// AndsImm - ands rd, rn, #imm (the flag-setting and).
func (p *Program) AndsImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.AndsImm(rd, rn, imm)
	return p.instrLine("ands", i, err, pos)
}

// AndsShift - ands rd, rn, rm[, shift #imm].
func (p *Program) AndsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("ands", err)
	}

	i, err := p.b.AndsShift(rd, rn, rm, v, sh)
	return p.instrLine("ands", i, err, pos)
}

// BicShift - bic rd, rn, rm[, shift #imm] (and-not).
func (p *Program) BicShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("bic", err)
	}

	i, err := p.b.BicShift(rd, rn, rm, v, sh)
	return p.instrLine("bic", i, err, pos)
}

// BicsShift - bics rd, rn, rm[, shift #imm] (the flag-setting bic).
func (p *Program) BicsShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("bics", err)
	}

	i, err := p.b.BicsShift(rd, rn, rm, v, sh)
	return p.instrLine("bics", i, err, pos)
}

// OrrImm - orr rd, rn, #imm (the bitmask immediate).
func (p *Program) OrrImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.OrrImm(rd, rn, imm)
	return p.instrLine("orr", i, err, pos)
}

// OrrShift - orr rd, rn, rm[, shift #imm].
func (p *Program) OrrShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("orr", err)
	}

	i, err := p.b.OrrShift(rd, rn, rm, v, sh)
	return p.instrLine("orr", i, err, pos)
}

// OrnShift - orn rd, rn, rm[, shift #imm] (or-not).
func (p *Program) OrnShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("orn", err)
	}

	i, err := p.b.OrnShift(rd, rn, rm, v, sh)
	return p.instrLine("orn", i, err, pos)
}

// EorImm - eor rd, rn, #imm (the bitmask immediate).
func (p *Program) EorImm(rd, rn arch.Reg, imm uint64) *Program {
	pos := p.pos()
	i, err := p.b.EorImm(rd, rn, imm)
	return p.instrLine("eor", i, err, pos)
}

// EorShift - eor rd, rn, rm[, shift #imm].
func (p *Program) EorShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("eor", err)
	}

	i, err := p.b.EorShift(rd, rn, rm, v, sh)
	return p.instrLine("eor", i, err, pos)
}

// EonShift - eon rd, rn, rm[, shift #imm] (eor-not).
func (p *Program) EonShift(rd, rn, rm arch.Reg, imm int64, sh arch.Shift) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(imm)
	if err != nil {
		return p.fail("eon", err)
	}

	i, err := p.b.EonShift(rd, rn, rm, v, sh)
	return p.instrLine("eon", i, err, pos)
}

// Ccmp - ccmp rn, rm, #nzcv, cond.
func (p *Program) Ccmp(rn, rm arch.Reg, nzcv uint32, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Ccmp(rn, rm, nzcv, cond)
	return p.instrLine("ccmp", i, err, pos)
}

// Csel - csel rd, rn, rm, cond.
func (p *Program) Csel(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csel(rd, rn, rm, cond)
	return p.instrLine("csel", i, err, pos)
}

// Csinc - csinc rd, rn, rm, cond.
func (p *Program) Csinc(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csinc(rd, rn, rm, cond)
	return p.instrLine("csinc", i, err, pos)
}

// Csinv - csinv rd, rn, rm, cond.
func (p *Program) Csinv(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csinv(rd, rn, rm, cond)
	return p.instrLine("csinv", i, err, pos)
}

// Csneg - csneg rd, rn, rm, cond.
func (p *Program) Csneg(rd, rn, rm arch.Reg, cond string) *Program {
	pos := p.pos()
	i, err := p.b.Csneg(rd, rn, rm, cond)
	return p.instrLine("csneg", i, err, pos)
}

// --- shifts / bitfield ----------------------------------------------------------

// LslReg - lsl rd, rn, rm (the register shift form).
func (p *Program) LslReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.LslReg(rd, rn, rm)
	return p.instrLine("lsl", i, err, pos)
}

// LsrReg - lsr rd, rn, rm (the register shift form).
func (p *Program) LsrReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.LsrReg(rd, rn, rm)
	return p.instrLine("lsr", i, err, pos)
}

// AsrReg - asr rd, rn, rm (the register shift form).
func (p *Program) AsrReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.AsrReg(rd, rn, rm)
	return p.instrLine("asr", i, err, pos)
}

// RorReg - ror rd, rn, rm (the register shift form).
func (p *Program) RorReg(rd, rn, rm arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.RorReg(rd, rn, rm)
	return p.instrLine("ror", i, err, pos)
}

// Bfm - bfm rd, rn, #immr, #imms.
func (p *Program) Bfm(rd, rn arch.Reg, immr, imms uint32) *Program {
	pos := p.pos()
	i, err := p.b.Bfm(rd, rn, immr, imms)
	return p.instrLine("bfm", i, err, pos)
}

// Sbfm - sbfm rd, rn, #immr, #imms.
func (p *Program) Sbfm(rd, rn arch.Reg, immr, imms uint32) *Program {
	pos := p.pos()
	i, err := p.b.Sbfm(rd, rn, immr, imms)
	return p.instrLine("sbfm", i, err, pos)
}

// Ubfm - ubfm rd, rn, #immr, #imms.
func (p *Program) Ubfm(rd, rn arch.Reg, immr, imms uint32) *Program {
	pos := p.pos()
	i, err := p.b.Ubfm(rd, rn, immr, imms)
	return p.instrLine("ubfm", i, err, pos)
}

// Extr - extr rd, rn, rm, #lsb.
func (p *Program) Extr(rd, rn, rm arch.Reg, lsb int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm6(lsb)
	if err != nil {
		return p.fail("extr", err)
	}

	i, err := p.b.Extr(rd, rn, rm, v)
	return p.instrLine("extr", i, err, pos)
}

// Cls - cls rd, rn.
func (p *Program) Cls(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Cls(rd, rn)
	return p.instrLine("cls", i, err, pos)
}

// Clz - clz rd, rn.
func (p *Program) Clz(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Clz(rd, rn)
	return p.instrLine("clz", i, err, pos)
}

// Rbit - rbit rd, rn.
func (p *Program) Rbit(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rbit(rd, rn)
	return p.instrLine("rbit", i, err, pos)
}

// Rev - rev rd, rn.
func (p *Program) Rev(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rev(rd, rn)
	return p.instrLine("rev", i, err, pos)
}

// Rev16 - rev16 rd, rn.
func (p *Program) Rev16(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rev16(rd, rn)
	return p.instrLine("rev16", i, err, pos)
}

// Rev32 - rev32 rd, rn.
func (p *Program) Rev32(rd, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Rev32(rd, rn)
	return p.instrLine("rev32", i, err, pos)
}

// --- control transfers ----------------------------------------------------------

// Br - br rn (branch to register).
func (p *Program) Br(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Br(rn)
	return p.instrLine("br", i, err, pos)
}

// Blr - blr rn (branch with link to register).
func (p *Program) Blr(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Blr(rn)
	return p.instrLine("blr", i, err, pos)
}

// Ret - ret rn.
func (p *Program) Ret(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ret(rn)
	return p.instrLine("ret", i, err, pos)
}

// --- loads ----------------------------------------------------------------------

// Ldr - ldr rt, [rn, #off].
func (p *Program) Ldr(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldr(rt, rn, arch.Off(off))
	return p.instrLine("ldr", i, err, pos)
}

// Ldrb - ldrb rt, [rn, #off].
func (p *Program) Ldrb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrb(rt, rn, arch.Off(off))
	return p.instrLine("ldrb", i, err, pos)
}

// Ldrh - ldrh rt, [rn, #off].
func (p *Program) Ldrh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrh(rt, rn, arch.Off(off))
	return p.instrLine("ldrh", i, err, pos)
}

// Ldrsb - ldrsb rt, [rn, #off].
func (p *Program) Ldrsb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrsb(rt, rn, arch.Off(off))
	return p.instrLine("ldrsb", i, err, pos)
}

// Ldrsh - ldrsh rt, [rn, #off].
func (p *Program) Ldrsh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrsh(rt, rn, arch.Off(off))
	return p.instrLine("ldrsh", i, err, pos)
}

// Ldrsw - ldrsw rt, [rn, #off].
func (p *Program) Ldrsw(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldrsw(rt, rn, arch.Off(off))
	return p.instrLine("ldrsw", i, err, pos)
}

// Ldur - ldur rt, [rn, #off] (unscaled).
func (p *Program) Ldur(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldur(rt, rn, arch.Off(off))
	return p.instrLine("ldur", i, err, pos)
}

// Ldurb - ldurb rt, [rn, #off] (unscaled).
func (p *Program) Ldurb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldurb(rt, rn, arch.Off(off))
	return p.instrLine("ldurb", i, err, pos)
}

// Ldurh - ldurh rt, [rn, #off] (unscaled).
func (p *Program) Ldurh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldurh(rt, rn, arch.Off(off))
	return p.instrLine("ldurh", i, err, pos)
}

// Ldp - ldp rt, rt2, [rn, #off].
func (p *Program) Ldp(rt, rt2, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldp(rt, rt2, rn, arch.Off(off))
	return p.instrLine("ldp", i, err, pos)
}

// Ldpsw - ldpsw rt, rt2, [rn, #off].
func (p *Program) Ldpsw(rt, rt2, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Ldpsw(rt, rt2, rn, arch.Off(off))
	return p.instrLine("ldpsw", i, err, pos)
}

// --- stores ---------------------------------------------------------------------

// Str - str rt, [rn, #off].
func (p *Program) Str(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Str(rt, rn, arch.Off(off))
	return p.instrLine("str", i, err, pos)
}

// Strb - strb rt, [rn, #off].
func (p *Program) Strb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Strb(rt, rn, arch.Off(off))
	return p.instrLine("strb", i, err, pos)
}

// Strh - strh rt, [rn, #off].
func (p *Program) Strh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Strh(rt, rn, arch.Off(off))
	return p.instrLine("strh", i, err, pos)
}

// Stur - stur rt, [rn, #off] (unscaled).
func (p *Program) Stur(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Stur(rt, rn, arch.Off(off))
	return p.instrLine("stur", i, err, pos)
}

// Sturb - sturb rt, [rn, #off] (unscaled).
func (p *Program) Sturb(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Sturb(rt, rn, arch.Off(off))
	return p.instrLine("sturb", i, err, pos)
}

// Sturh - sturh rt, [rn, #off] (unscaled).
func (p *Program) Sturh(rt, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Sturh(rt, rn, arch.Off(off))
	return p.instrLine("sturh", i, err, pos)
}

// Stp - stp rt, rt2, [rn, #off].
func (p *Program) Stp(rt, rt2, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.Stp(rt, rt2, rn, arch.Off(off))
	return p.instrLine("stp", i, err, pos)
}

// --- atomics --------------------------------------------------------------------

// Ldar - ldar rt, [rn].
func (p *Program) Ldar(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldar(rt, rn)
	return p.instrLine("ldar", i, err, pos)
}

// Ldarb - ldarb rt, [rn].
func (p *Program) Ldarb(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldarb(rt, rn)
	return p.instrLine("ldarb", i, err, pos)
}

// Ldaxr - ldaxr rt, [rn].
func (p *Program) Ldaxr(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldaxr(rt, rn)
	return p.instrLine("ldaxr", i, err, pos)
}

// Ldaxrb - ldaxrb rt, [rn].
func (p *Program) Ldaxrb(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ldaxrb(rt, rn)
	return p.instrLine("ldaxrb", i, err, pos)
}

// Stlr - stlr rt, [rn].
func (p *Program) Stlr(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlr(rt, rn)
	return p.instrLine("stlr", i, err, pos)
}

// Stlrb - stlrb rt, [rn].
func (p *Program) Stlrb(rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlrb(rt, rn)
	return p.instrLine("stlrb", i, err, pos)
}

// Stlxr - stlxr rs, rt, [rn].
func (p *Program) Stlxr(rs, rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlxr(rs, rt, rn)
	return p.instrLine("stlxr", i, err, pos)
}

// Stlxrb - stlxrb rs, rt, [rn].
func (p *Program) Stlxrb(rs, rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stlxrb(rs, rt, rn)
	return p.instrLine("stlxrb", i, err, pos)
}

// Stxrb - stxrb rs, rt, [rn].
func (p *Program) Stxrb(rs, rt, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Stxrb(rs, rt, rn)
	return p.instrLine("stxrb", i, err, pos)
}

// --- floating point ------------------------------------------------------------

// Fadd - fadd fd, fn, fm (double/single by the operand kind).
func (p *Program) Fadd(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fadd(rd, rn, rm)
	return p.instrLine("fadd", i, err, pos)
}

// Fsub - fsub fd, fn, fm.
func (p *Program) Fsub(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fsub(rd, rn, rm)
	return p.instrLine("fsub", i, err, pos)
}

// Fmul - fmul fd, fn, fm.
func (p *Program) Fmul(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmul(rd, rn, rm)
	return p.instrLine("fmul", i, err, pos)
}

// Fdiv - fdiv fd, fn, fm.
func (p *Program) Fdiv(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fdiv(rd, rn, rm)
	return p.instrLine("fdiv", i, err, pos)
}

// Fmax - fmax fd, fn, fm.
func (p *Program) Fmax(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmax(rd, rn, rm)
	return p.instrLine("fmax", i, err, pos)
}

// Fmin - fmin fd, fn, fm.
func (p *Program) Fmin(rd, rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmin(rd, rn, rm)
	return p.instrLine("fmin", i, err, pos)
}

// Fneg - fneg fd, fn.
func (p *Program) Fneg(rd, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fneg(rd, rn)
	return p.instrLine("fneg", i, err, pos)
}

// Fmov - fmov fd, fn (the register form between two FP registers).
func (p *Program) Fmov(rd, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmov(rd, rn)
	return p.instrLine("fmov", i, err, pos)
}

// FmovFromGpr - fmov fd, xn | fmov sn, wn (bits from the integer file;
// fmov d0, xzr is the FP zero).
func (p *Program) FmovFromGpr(rd arch.FReg, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.FmovFromGpr(rd, rn)
	return p.instrLine("fmov", i, err, pos)
}

// FmovToGpr - fmov xn, fd | fmov wn, sn (bits to the integer file).
func (p *Program) FmovToGpr(rd arch.Reg, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.FmovToGpr(rd, rn)
	return p.instrLine("fmov", i, err, pos)
}

// FmovImm - fmov fd, #imm (the VFP imm8 form).
func (p *Program) FmovImm(rd arch.FReg, imm float64) *Program {
	pos := p.pos()
	i, err := p.b.FmovImm(rd, imm)
	return p.instrLine("fmov", i, err, pos)
}

// Fcvt - fcvt fd, fn (the s<->d width conversion).
func (p *Program) Fcvt(rd, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcvt(rd, rn)
	return p.instrLine("fcvt", i, err, pos)
}

// Scvtf - scvtf fd, wn|xn (signed integer to FP).
func (p *Program) Scvtf(rd arch.FReg, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Scvtf(rd, rn)
	return p.instrLine("scvtf", i, err, pos)
}

// Ucvtf - ucvtf fd, wn|xn (unsigned integer to FP).
func (p *Program) Ucvtf(rd arch.FReg, rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Ucvtf(rd, rn)
	return p.instrLine("ucvtf", i, err, pos)
}

// Fcvtzs - fcvtzs wd|xd, fn (FP to signed integer).
func (p *Program) Fcvtzs(rd arch.Reg, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcvtzs(rd, rn)
	return p.instrLine("fcvtzs", i, err, pos)
}

// Fcvtzu - fcvtzu wd|xd, fn (FP to unsigned integer).
func (p *Program) Fcvtzu(rd arch.Reg, rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcvtzu(rd, rn)
	return p.instrLine("fcvtzu", i, err, pos)
}

// Fmadd - fmadd fd, fn, fm, fa (fd = fa + fn*fm).
func (p *Program) Fmadd(rd, rn, rm, ra arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fmadd(rd, rn, rm, ra)
	return p.instrLine("fmadd", i, err, pos)
}

// Fnmsub - fnmsub fd, fn, fm, fa (fd = -(fn*fm - fa)).
func (p *Program) Fnmsub(rd, rn, rm, ra arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fnmsub(rd, rn, rm, ra)
	return p.instrLine("fnmsub", i, err, pos)
}

// Fcmp - fcmp fn, fm (sets the FP flags).
func (p *Program) Fcmp(rn, rm arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.Fcmp(rn, rm)
	return p.instrLine("fcmp", i, err, pos)
}

// FcmpZero - fcmp fn, #0.0.
func (p *Program) FcmpZero(rn arch.FReg) *Program {
	pos := p.pos()
	i, err := p.b.FcmpZero(rn)
	return p.instrLine("fcmp", i, err, pos)
}

// LdrF - ldr st|dt, [rn, #off] (the FP/SIMD register form).
func (p *Program) LdrF(rt arch.FReg, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.LdrF(rt, rn, arch.Off(off))
	return p.instrLine("ldr", i, err, pos)
}

// StrF - str st|dt, [rn, #off] (the FP/SIMD register form).
func (p *Program) StrF(rt arch.FReg, rn arch.Reg, off int64) *Program {
	pos := p.pos()
	i, err := p.b.StrF(rt, rn, arch.Off(off))
	return p.instrLine("str", i, err, pos)
}

// --- system ---------------------------------------------------------------------

// Nop - nop.
func (p *Program) Nop() *Program {
	pos := p.pos()
	return p.instrLine("nop", p.b.Nop(), nil, pos)
}

// Brk - brk #imm.
func (p *Program) Brk(imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("brk", err)
	}

	return p.instrLine("brk", p.b.Brk(v), nil, pos)
}

// Mrs - mrs rd, sysreg.
func (p *Program) Mrs(rd arch.Reg, sysreg string) *Program {
	pos := p.pos()
	i, err := p.b.Mrs(rd, sysreg)
	return p.instrLine("mrs", i, err, pos)
}

// Msr - msr sysreg, rt.
func (p *Program) Msr(sysreg string, rt arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Msr(sysreg, rt)
	return p.instrLine("msr", i, err, pos)
}

// Prfm - prfm [rn].
func (p *Program) Prfm(rn arch.Reg) *Program {
	pos := p.pos()
	i, err := p.b.Prfm(rn)
	return p.instrLine("prfm", i, err, pos)
}

// Build - materialize the program; deferred construction errors are
// returned alongside (an empty slice means the program is well-formed).
func (p *Program) Build() (*Binary, []error) {
	return &Binary{Entry: p.entry, lines: p.lines}, p.errs
}

// --- internals ---------------------------------------------------------------

func (p *Program) instrLine(src string, i arch.Instr, err error, pos prog.Pos) *Program {
	if err != nil {
		return p.fail(src, err)
	}

	if p.stream != 0 {
		return p.fail(src, errors.New("an instruction in the data stream"))
	}

	p.lines = append(p.lines, newInstrLine(i, src, pos))
	return p
}

func (p *Program) branchLine(
	src, label string,
	ctor func(target, pc uint64) (arch.Instr, error),
	pos prog.Pos,
) *Program {
	if p.stream != 0 {
		return p.fail(src, errors.New("an instruction in the data stream"))
	}

	p.lines = append(p.lines, newBranchLine(src, label, ctor, pos))
	return p
}

func (p *Program) fail(src string, err error) *Program {
	p.errs = append(p.errs, fmt.Errorf("%s: %w", src, err))
	return p
}
