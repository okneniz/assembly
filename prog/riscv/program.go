// Package riscv - programs written directly in Go: a chain Program over
// the arch builders, labels resolved at assembly time. The counterpart of
// an .s source file; immediate operands are raw ints, validation errors
// are deferred to Assemble. Every line encodes in its fixed 4-byte form
// (NoRVC): the program is laid out before encoding, and compressed forms
// would make line sizes depend on the pass - the syntax layer's policy
// for symbolic slots, applied to every line of a chain program. The Go
// source position of each line (the debugger's address ↔ source map)
// comes from the position resolver injected with WithPos - there is no
// built-in caller detection.
package riscv

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/riscv"
	"github.com/okneniz/assembly/prog"
)

// Program - a program being built: a sequence of lines (instructions,
// label-directed branches, la pairs, labels, data). Chain methods append
// one line each and return the program; nothing is encoded until
// Assemble.
type Program struct {
	lines []line
	entry string
	errs  []error
	pos   func() prog.Pos
	b     arch.Builder
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

// Label - define a label at the current position.
func (p *Program) Label(name string) *Program {
	p.lines = append(p.lines, newLabelLine(name, p.pos()))
	return p
}

// Entry - the label the program starts at.
func (p *Program) Entry(name string) *Program {
	p.entry = name
	return p
}

// Ascii - string data appended verbatim (no terminating zero).
func (p *Program) Ascii(s string) *Program {
	p.lines = append(p.lines, newDataLine([]byte(s), ".ascii", p.pos()))
	return p
}

// Bytes - raw data bytes.
func (p *Program) Bytes(b ...byte) *Program {
	p.lines = append(p.lines, newDataLine(b, ".byte", p.pos()))
	return p
}

// --- label-directed lines -------------------------------------------------------

// Jal - jump with link to a label (the return address lands in ra; the
// unconditional jump is Jal with the zero register).
func (p *Program) Jal(rd arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("jal", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Jal(rd, int64(t)-int64(pc)), nil
	}, pos)
}

// Beq - branch to a label when the registers are equal.
func (p *Program) Beq(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("beq", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Beq(rs1, rs2, int64(t)-int64(pc)), nil
	}, pos)
}

// Bne - branch to a label when the registers are not equal.
func (p *Program) Bne(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bne", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bne(rs1, rs2, int64(t)-int64(pc)), nil
	}, pos)
}

// Bge - branch to a label when rs1 >= rs2 (signed).
func (p *Program) Bge(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bge", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bge(rs1, rs2, int64(t)-int64(pc)), nil
	}, pos)
}

// Bgeu - branch to a label when rs1 >= rs2 (unsigned).
func (p *Program) Bgeu(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bgeu", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bgeu(rs1, rs2, int64(t)-int64(pc)), nil
	}, pos)
}

// Blt - branch to a label when rs1 < rs2 (signed).
func (p *Program) Blt(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("blt", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Blt(rs1, rs2, int64(t)-int64(pc)), nil
	}, pos)
}

// Bltu - branch to a label when rs1 < rs2 (unsigned).
func (p *Program) Bltu(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bltu", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bltu(rs1, rs2, int64(t)-int64(pc)), nil
	}, pos)
}

// La - load the address of a label into the register: the auipc+addi
// pair (a fixed 8 bytes; the split is computed against the pair's own
// address, exactly as the text-path pseudo).
func (p *Program) La(rd arch.Reg, label string) *Program {
	pos := p.pos()
	p.lines = append(p.lines, newLaLine(label, func(t, pc uint64) ([]arch.Instr, error) {
		return laPair(p.b, rd, int64(t), int64(pc))
	}, pos))
	return p
}

// laPair - the evaluated la encoding: auipc (hi) + addi (lo).
func laPair(b arch.Builder, rd arch.Reg, target, pc int64) ([]arch.Instr, error) {
	hi, lo := arch.PcrelHiLo(target - pc)
	hi20, err := b.Imm20(hi & 0xfffff)
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	lo12, err := b.Imm12(lo)
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	return []arch.Instr{
		b.Auipc(rd, hi20),
		b.Addi(rd, rd, lo12),
	}, nil
}

// --- upper immediates -----------------------------------------------------------

// Lui - lui rd, imm (bits 31:12 of the constant: 0x10000000 is Lui(A0,
// 0x10000); values with a low half come from Lui+Addi pairs).
func (p *Program) Lui(rd arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm20(imm)
	if err != nil {
		return p.fail("lui", err)
	}

	return p.instrLine("lui", p.b.Lui(rd, v), nil, pos)
}

// Auipc - auipc rd, imm (pc + the 20-bit upper immediate; the label
// form is the La pair).
func (p *Program) Auipc(rd arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm20(imm)
	if err != nil {
		return p.fail("auipc", err)
	}

	return p.instrLine("auipc", p.b.Auipc(rd, v), nil, pos)
}

// --- arithmetic -----------------------------------------------------------------

// Add - add rd, rs1, rs2.
func (p *Program) Add(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("add", p.b.Add(rd, rs1, rs2), nil, pos)
}

// Addw - addw rd, rs1, rs2 (32-bit result, sign-extended).
func (p *Program) Addw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("addw", p.b.Addw(rd, rs1, rs2), nil, pos)
}

// Addi - addi rd, rs1, imm (rs1 = zero prints as li, imm = 0 - mv).
func (p *Program) Addi(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("addi", err)
	}

	return p.instrLine("addi", p.b.Addi(rd, rs1, v), nil, pos)
}

// Addiw - addiw rd, rs1, imm (32-bit, sign-extended).
func (p *Program) Addiw(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("addiw", err)
	}

	return p.instrLine("addiw", p.b.Addiw(rd, rs1, v), nil, pos)
}

// Sub - sub rd, rs1, rs2.
func (p *Program) Sub(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("sub", p.b.Sub(rd, rs1, rs2), nil, pos)
}

// Subw - subw rd, rs1, rs2 (32-bit result, sign-extended).
func (p *Program) Subw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("subw", p.b.Subw(rd, rs1, rs2), nil, pos)
}

// Mul - mul rd, rs1, rs2 (the low 64 bits of the product).
func (p *Program) Mul(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("mul", p.b.Mul(rd, rs1, rs2), nil, pos)
}

// Mulh - mulh rd, rs1, rs2 (the high 64 bits, signed).
func (p *Program) Mulh(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("mulh", p.b.Mulh(rd, rs1, rs2), nil, pos)
}

// Mulhsu - mulhsu rd, rs1, rs2 (the high 64 bits, rs1 signed).
func (p *Program) Mulhsu(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("mulhsu", p.b.Mulhsu(rd, rs1, rs2), nil, pos)
}

// Mulhu - mulhu rd, rs1, rs2 (the high 64 bits, unsigned).
func (p *Program) Mulhu(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("mulhu", p.b.Mulhu(rd, rs1, rs2), nil, pos)
}

// Mulw - mulw rd, rs1, rs2 (32-bit product, sign-extended).
func (p *Program) Mulw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("mulw", p.b.Mulw(rd, rs1, rs2), nil, pos)
}

// Div - div rd, rs1, rs2 (signed).
func (p *Program) Div(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("div", p.b.Div(rd, rs1, rs2), nil, pos)
}

// Divu - divu rd, rs1, rs2 (unsigned).
func (p *Program) Divu(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("divu", p.b.Divu(rd, rs1, rs2), nil, pos)
}

// Divw - divw rd, rs1, rs2 (32-bit, signed).
func (p *Program) Divw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("divw", p.b.Divw(rd, rs1, rs2), nil, pos)
}

// Divuw - divuw rd, rs1, rs2 (32-bit, unsigned).
func (p *Program) Divuw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("divuw", p.b.Divuw(rd, rs1, rs2), nil, pos)
}

// Rem - rem rd, rs1, rs2 (signed remainder).
func (p *Program) Rem(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("rem", p.b.Rem(rd, rs1, rs2), nil, pos)
}

// Remu - remu rd, rs1, rs2 (unsigned remainder).
func (p *Program) Remu(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("remu", p.b.Remu(rd, rs1, rs2), nil, pos)
}

// Remw - remw rd, rs1, rs2 (32-bit, signed).
func (p *Program) Remw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("remw", p.b.Remw(rd, rs1, rs2), nil, pos)
}

// Remuw - remuw rd, rs1, rs2 (32-bit, unsigned).
func (p *Program) Remuw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("remuw", p.b.Remuw(rd, rs1, rs2), nil, pos)
}

// Slt - slt rd, rs1, rs2 (1 when rs1 < rs2, signed).
func (p *Program) Slt(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("slt", p.b.Slt(rd, rs1, rs2), nil, pos)
}

// Sltu - sltu rd, rs1, rs2 (1 when rs1 < rs2, unsigned).
func (p *Program) Sltu(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("sltu", p.b.Sltu(rd, rs1, rs2), nil, pos)
}

// Slti - slti rd, rs1, imm (signed compare).
func (p *Program) Slti(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("slti", err)
	}

	return p.instrLine("slti", p.b.Slti(rd, rs1, v), nil, pos)
}

// Sltiu - sltiu rd, rs1, imm (unsigned compare).
func (p *Program) Sltiu(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("sltiu", err)
	}

	return p.instrLine("sltiu", p.b.Sltiu(rd, rs1, v), nil, pos)
}

// --- logical --------------------------------------------------------------------

// And - and rd, rs1, rs2.
func (p *Program) And(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("and", p.b.And(rd, rs1, rs2), nil, pos)
}

// Andi - andi rd, rs1, imm.
func (p *Program) Andi(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("andi", err)
	}

	return p.instrLine("andi", p.b.Andi(rd, rs1, v), nil, pos)
}

// Or - or rd, rs1, rs2.
func (p *Program) Or(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("or", p.b.Or(rd, rs1, rs2), nil, pos)
}

// Ori - ori rd, rs1, imm.
func (p *Program) Ori(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("ori", err)
	}

	return p.instrLine("ori", p.b.Ori(rd, rs1, v), nil, pos)
}

// Xor - xor rd, rs1, rs2.
func (p *Program) Xor(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("xor", p.b.Xor(rd, rs1, rs2), nil, pos)
}

// Xori - xori rd, rs1, imm.
func (p *Program) Xori(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("xori", err)
	}

	return p.instrLine("xori", p.b.Xori(rd, rs1, v), nil, pos)
}

// --- shifts ---------------------------------------------------------------------

// Sll - sll rd, rs1, rs2 (shift left by rs2 & 63).
func (p *Program) Sll(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("sll", p.b.Sll(rd, rs1, rs2), nil, pos)
}

// Sllw - sllw rd, rs1, rs2 (32-bit, shift amount & 31).
func (p *Program) Sllw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("sllw", p.b.Sllw(rd, rs1, rs2), nil, pos)
}

// Slli - slli rd, rs1, shamt.
func (p *Program) Slli(rd, rs1 arch.Reg, shamt int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(shamt)
	if err != nil {
		return p.fail("slli", err)
	}

	return p.instrLine("slli", p.b.Slli(rd, rs1, v), nil, pos)
}

// Slliw - slliw rd, rs1, shamt (32-bit).
func (p *Program) Slliw(rd, rs1 arch.Reg, shamt int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(shamt)
	if err != nil {
		return p.fail("slliw", err)
	}

	return p.instrLine("slliw", p.b.Slliw(rd, rs1, v), nil, pos)
}

// Srl - srl rd, rs1, rs2 (shift right logical).
func (p *Program) Srl(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("srl", p.b.Srl(rd, rs1, rs2), nil, pos)
}

// Srlw - srlw rd, rs1, rs2 (32-bit).
func (p *Program) Srlw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("srlw", p.b.Srlw(rd, rs1, rs2), nil, pos)
}

// Srli - srli rd, rs1, shamt.
func (p *Program) Srli(rd, rs1 arch.Reg, shamt int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(shamt)
	if err != nil {
		return p.fail("srli", err)
	}

	return p.instrLine("srli", p.b.Srli(rd, rs1, v), nil, pos)
}

// Srliw - srliw rd, rs1, shamt (32-bit).
func (p *Program) Srliw(rd, rs1 arch.Reg, shamt int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(shamt)
	if err != nil {
		return p.fail("srliw", err)
	}

	return p.instrLine("srliw", p.b.Srliw(rd, rs1, v), nil, pos)
}

// Sra - sra rd, rs1, rs2 (shift right arithmetic).
func (p *Program) Sra(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("sra", p.b.Sra(rd, rs1, rs2), nil, pos)
}

// Sraw - sraw rd, rs1, rs2 (32-bit).
func (p *Program) Sraw(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("sraw", p.b.Sraw(rd, rs1, rs2), nil, pos)
}

// Srai - srai rd, rs1, shamt.
func (p *Program) Srai(rd, rs1 arch.Reg, shamt int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(shamt)
	if err != nil {
		return p.fail("srai", err)
	}

	return p.instrLine("srai", p.b.Srai(rd, rs1, v), nil, pos)
}

// Sraiw - sraiw rd, rs1, shamt (32-bit).
func (p *Program) Sraiw(rd, rs1 arch.Reg, shamt int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(shamt)
	if err != nil {
		return p.fail("sraiw", err)
	}

	return p.instrLine("sraiw", p.b.Sraiw(rd, rs1, v), nil, pos)
}

// --- jumps ----------------------------------------------------------------------

// Jalr - jalr rd, off(rs1) (the indirect jump; the plain "jalr rs"
// pseudo form is Jalr with rd = ra and off = 0).
func (p *Program) Jalr(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("jalr", err)
	}

	return p.instrLine("jalr", p.b.Jalr(rd, rs1, v), nil, pos)
}

// --- loads ----------------------------------------------------------------------

// Lb - lb rd, off(rs1) (a sign-extended byte).
func (p *Program) Lb(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lb", err)
	}

	return p.instrLine("lb", p.b.Lb(rd, rs1, v), nil, pos)
}

// Lbu - lbu rd, off(rs1) (a zero-extended byte).
func (p *Program) Lbu(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lbu", err)
	}

	return p.instrLine("lbu", p.b.Lbu(rd, rs1, v), nil, pos)
}

// Lh - lh rd, off(rs1) (a sign-extended halfword).
func (p *Program) Lh(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lh", err)
	}

	return p.instrLine("lh", p.b.Lh(rd, rs1, v), nil, pos)
}

// Lhu - lhu rd, off(rs1) (a zero-extended halfword).
func (p *Program) Lhu(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lhu", err)
	}

	return p.instrLine("lhu", p.b.Lhu(rd, rs1, v), nil, pos)
}

// Lw - lw rd, off(rs1) (a sign-extended word).
func (p *Program) Lw(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lw", err)
	}

	return p.instrLine("lw", p.b.Lw(rd, rs1, v), nil, pos)
}

// Lwu - lwu rd, off(rs1) (a zero-extended word).
func (p *Program) Lwu(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lwu", err)
	}

	return p.instrLine("lwu", p.b.Lwu(rd, rs1, v), nil, pos)
}

// Ld - ld rd, off(rs1).
func (p *Program) Ld(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("ld", err)
	}

	return p.instrLine("ld", p.b.Ld(rd, rs1, v), nil, pos)
}

// Flw - flw fd, off(rs1) (a 32-bit float; the FP registers are Reg
// numbers 0..31, printed via the FP ABI table).
func (p *Program) Flw(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("flw", err)
	}

	return p.instrLine("flw", p.b.Flw(rd, rs1, v), nil, pos)
}

// Fld - fld fd, off(rs1) (a 64-bit float).
func (p *Program) Fld(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("fld", err)
	}

	return p.instrLine("fld", p.b.Fld(rd, rs1, v), nil, pos)
}

// --- stores ---------------------------------------------------------------------

// Sb - sb rs2, off(rs1) (the low byte of rs2).
func (p *Program) Sb(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("sb", err)
	}

	return p.instrLine("sb", p.b.Sb(rs2, rs1, v), nil, pos)
}

// Sh - sh rs2, off(rs1) (the low halfword).
func (p *Program) Sh(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("sh", err)
	}

	return p.instrLine("sh", p.b.Sh(rs2, rs1, v), nil, pos)
}

// Sw - sw rs2, off(rs1) (the low word).
func (p *Program) Sw(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("sw", err)
	}

	return p.instrLine("sw", p.b.Sw(rs2, rs1, v), nil, pos)
}

// Sd - sd rs2, off(rs1).
func (p *Program) Sd(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("sd", err)
	}

	return p.instrLine("sd", p.b.Sd(rs2, rs1, v), nil, pos)
}

// Fsw - fsw fs2, off(rs1) (a 32-bit float).
func (p *Program) Fsw(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("fsw", err)
	}

	return p.instrLine("fsw", p.b.Fsw(rs2, rs1, v), nil, pos)
}

// Fsd - fsd fs2, off(rs1) (a 64-bit float).
func (p *Program) Fsd(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("fsd", err)
	}

	return p.instrLine("fsd", p.b.Fsd(rs2, rs1, v), nil, pos)
}

// --- atomics --------------------------------------------------------------------

// AmoswapW - amoswap.w rd, rs2, (rs1) (32-bit swap).
func (p *Program) AmoswapW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoswap.w", p.b.AmoswapW(rd, rs1, rs2), nil, pos)
}

// AmoswapD - amoswap.d rd, rs2, (rs1) (64-bit swap).
func (p *Program) AmoswapD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoswap.d", p.b.AmoswapD(rd, rs1, rs2), nil, pos)
}

// AmoaddW - amoadd.w rd, rs2, (rs1).
func (p *Program) AmoaddW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoadd.w", p.b.AmoaddW(rd, rs1, rs2), nil, pos)
}

// AmoaddD - amoadd.d rd, rs2, (rs1).
func (p *Program) AmoaddD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoadd.d", p.b.AmoaddD(rd, rs1, rs2), nil, pos)
}

// AmoxorW - amoxor.w rd, rs2, (rs1).
func (p *Program) AmoxorW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoxor.w", p.b.AmoxorW(rd, rs1, rs2), nil, pos)
}

// AmoxorD - amoxor.d rd, rs2, (rs1).
func (p *Program) AmoxorD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoxor.d", p.b.AmoxorD(rd, rs1, rs2), nil, pos)
}

// AmoandW - amoand.w rd, rs2, (rs1).
func (p *Program) AmoandW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoand.w", p.b.AmoandW(rd, rs1, rs2), nil, pos)
}

// AmoandD - amoand.d rd, rs2, (rs1).
func (p *Program) AmoandD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoand.d", p.b.AmoandD(rd, rs1, rs2), nil, pos)
}

// AmoorW - amoor.w rd, rs2, (rs1).
func (p *Program) AmoorW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoor.w", p.b.AmoorW(rd, rs1, rs2), nil, pos)
}

// AmoorD - amoor.d rd, rs2, (rs1).
func (p *Program) AmoorD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amoor.d", p.b.AmoorD(rd, rs1, rs2), nil, pos)
}

// AmominW - amomin.w rd, rs2, (rs1) (signed minimum).
func (p *Program) AmominW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amomin.w", p.b.AmominW(rd, rs1, rs2), nil, pos)
}

// AmominD - amomin.d rd, rs2, (rs1) (signed minimum).
func (p *Program) AmominD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amomin.d", p.b.AmominD(rd, rs1, rs2), nil, pos)
}

// AmominuW - amominu.w rd, rs2, (rs1) (unsigned minimum).
func (p *Program) AmominuW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amominu.w", p.b.AmominuW(rd, rs1, rs2), nil, pos)
}

// AmominuD - amominu.d rd, rs2, (rs1) (unsigned minimum).
func (p *Program) AmominuD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amominu.d", p.b.AmominuD(rd, rs1, rs2), nil, pos)
}

// AmomaxW - amomax.w rd, rs2, (rs1) (signed maximum).
func (p *Program) AmomaxW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amomax.w", p.b.AmomaxW(rd, rs1, rs2), nil, pos)
}

// AmomaxD - amomax.d rd, rs2, (rs1) (signed maximum).
func (p *Program) AmomaxD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amomax.d", p.b.AmomaxD(rd, rs1, rs2), nil, pos)
}

// AmomaxuW - amomaxu.w rd, rs2, (rs1) (unsigned maximum).
func (p *Program) AmomaxuW(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amomaxu.w", p.b.AmomaxuW(rd, rs1, rs2), nil, pos)
}

// AmomaxuD - amomaxu.d rd, rs2, (rs1) (unsigned maximum).
func (p *Program) AmomaxuD(rd, rs1, rs2 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("amomaxu.d", p.b.AmomaxuD(rd, rs1, rs2), nil, pos)
}

// --- CSR ------------------------------------------------------------------------

// Csrrw - csrrw rd, csr, rs1 (atomic read/write).
func (p *Program) Csrrw(rd arch.Reg, csr uint16, rs1 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("csrrw", p.b.Csrrw(rd, csr, rs1), nil, pos)
}

// Csrrs - csrrs rd, csr, rs1 (read and set bits).
func (p *Program) Csrrs(rd arch.Reg, csr uint16, rs1 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("csrrs", p.b.Csrrs(rd, csr, rs1), nil, pos)
}

// Csrrc - csrrc rd, csr, rs1 (read and clear bits).
func (p *Program) Csrrc(rd arch.Reg, csr uint16, rs1 arch.Reg) *Program {
	pos := p.pos()
	return p.instrLine("csrrc", p.b.Csrrc(rd, csr, rs1), nil, pos)
}

// Csrrwi - csrrwi rd, csr, zimm (immediate write).
func (p *Program) Csrrwi(rd arch.Reg, csr uint16, zimm uint8) *Program {
	pos := p.pos()
	return p.instrLine("csrrwi", p.b.Csrrwi(rd, csr, zimm), nil, pos)
}

// Csrrsi - csrrsi rd, csr, zimm (immediate set bits).
func (p *Program) Csrrsi(rd arch.Reg, csr uint16, zimm uint8) *Program {
	pos := p.pos()
	return p.instrLine("csrrsi", p.b.Csrrsi(rd, csr, zimm), nil, pos)
}

// Csrrci - csrrci rd, csr, zimm (immediate clear bits).
func (p *Program) Csrrci(rd arch.Reg, csr uint16, zimm uint8) *Program {
	pos := p.pos()
	return p.instrLine("csrrci", p.b.Csrrci(rd, csr, zimm), nil, pos)
}

// --- floating point -------------------------------------------------------------

// FaddS - fadd.s fd, fs1, fs2, rm (rm = the rounding mode 0..7).
func (p *Program) FaddS(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fadd.s", p.b.FaddS(rd, rs1, rs2, rm), nil, pos)
}

// FaddD - fadd.d fd, fs1, fs2, rm.
func (p *Program) FaddD(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fadd.d", p.b.FaddD(rd, rs1, rs2, rm), nil, pos)
}

// FsubS - fsub.s fd, fs1, fs2, rm.
func (p *Program) FsubS(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fsub.s", p.b.FsubS(rd, rs1, rs2, rm), nil, pos)
}

// FsubD - fsub.d fd, fs1, fs2, rm.
func (p *Program) FsubD(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fsub.d", p.b.FsubD(rd, rs1, rs2, rm), nil, pos)
}

// FmulS - fmul.s fd, fs1, fs2, rm.
func (p *Program) FmulS(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fmul.s", p.b.FmulS(rd, rs1, rs2, rm), nil, pos)
}

// FmulD - fmul.d fd, fs1, fs2, rm.
func (p *Program) FmulD(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fmul.d", p.b.FmulD(rd, rs1, rs2, rm), nil, pos)
}

// FdivS - fdiv.s fd, fs1, fs2, rm.
func (p *Program) FdivS(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fdiv.s", p.b.FdivS(rd, rs1, rs2, rm), nil, pos)
}

// FdivD - fdiv.d fd, fs1, fs2, rm.
func (p *Program) FdivD(rd, rs1, rs2 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fdiv.d", p.b.FdivD(rd, rs1, rs2, rm), nil, pos)
}

// FmaddS - fmadd.s fd, fs1, fs2, fs3, rm (fd = fs1*fs2 + fs3).
func (p *Program) FmaddS(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fmadd.s", p.b.FmaddS(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FmaddD - fmadd.d fd, fs1, fs2, fs3, rm.
func (p *Program) FmaddD(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fmadd.d", p.b.FmaddD(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FmsubS - fmsub.s fd, fs1, fs2, fs3, rm (fd = fs1*fs2 - fs3).
func (p *Program) FmsubS(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fmsub.s", p.b.FmsubS(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FmsubD - fmsub.d fd, fs1, fs2, fs3, rm.
func (p *Program) FmsubD(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fmsub.d", p.b.FmsubD(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FnmaddS - fnmadd.s fd, fs1, fs2, fs3, rm (fd = -fs1*fs2 + fs3).
func (p *Program) FnmaddS(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fnmadd.s", p.b.FnmaddS(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FnmaddD - fnmadd.d fd, fs1, fs2, fs3, rm.
func (p *Program) FnmaddD(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fnmadd.d", p.b.FnmaddD(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FnmsubS - fnmsub.s fd, fs1, fs2, fs3, rm (fd = -fs1*fs2 - fs3).
func (p *Program) FnmsubS(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fnmsub.s", p.b.FnmsubS(rd, rs1, rs2, rs3, rm), nil, pos)
}

// FnmsubD - fnmsub.d fd, fs1, fs2, fs3, rm.
func (p *Program) FnmsubD(rd, rs1, rs2, rs3 arch.Reg, rm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fnmsub.d", p.b.FnmsubD(rd, rs1, rs2, rs3, rm), nil, pos)
}

// --- system ---------------------------------------------------------------------

// Ecall - ecall (the syscall trap).
func (p *Program) Ecall() *Program {
	pos := p.pos()
	i, err := arch.BuildInstr("ecall", nil)
	return p.instrLine("ecall", i, err, pos)
}

// Fence - fence fm (the memory-ordering fence; fm = the predecessor/
// successor ordering bits, e.g. iorw, iorw).
func (p *Program) Fence(fm uint8) *Program {
	pos := p.pos()
	return p.instrLine("fence", p.b.Fence(fm), nil, pos)
}

// Build - materialize the program; deferred construction errors are
// returned alongside.
func (p *Program) Build() (*Binary, []error) {
	return &Binary{Entry: p.entry, lines: p.lines}, p.errs
}

// --- internals ---------------------------------------------------------------

func (p *Program) instrLine(src string, i arch.Instr, err error, pos prog.Pos) *Program {
	if err != nil {
		return p.fail(src, err)
	}

	p.lines = append(p.lines, newInstrLine(i, src, pos))
	return p
}

func (p *Program) branchLine(
	src, label string,
	ctor func(target, pc uint64) (arch.Instr, error),
	pos prog.Pos,
) *Program {
	p.lines = append(p.lines, newBranchLine(src, label, ctor, pos))
	return p
}

func (p *Program) fail(src string, err error) *Program {
	p.errs = append(p.errs, fmt.Errorf("%s: %w", src, err))
	return p
}
