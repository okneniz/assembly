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

// Instr - append an already computed instruction (the escape hatch when
// no chain twin exists yet).
func (p *Program) Instr(i arch.Instr, src string) *Program {
	p.lines = append(p.lines, newInstrLine(i, src, p.pos()))
	return p
}

// --- label-directed lines -------------------------------------------------------

// J - unconditional jump to a label (jal zero - prints as j).
func (p *Program) J(label string) *Program {
	pos := p.pos()
	return p.branchLine("j", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Jal(arch.Zero, int64(t)-int64(pc)), nil
	}, pos)
}

// Jal - jump with link to a label (the return address lands in ra).
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

// Bnez - branch to a label when the register is not zero.
func (p *Program) Bnez(rs arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bnez", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bne(rs, arch.Zero, int64(t)-int64(pc)), nil
	}, pos)
}

// Bgeu - branch to a label when rs1 >= rs2 (unsigned).
func (p *Program) Bgeu(rs1, rs2 arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bgeu", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bgeu(rs1, rs2, int64(t)-int64(pc)), nil
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

// --- instruction chain twins ----------------------------------------------------

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

// Addi - addi rd, rs1, imm (rs1 = zero prints as li, imm = 0 - mv).
func (p *Program) Addi(rd, rs1 arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("addi", err)
	}

	return p.instrLine("addi", p.b.Addi(rd, rs1, v), nil, pos)
}

// Lb - lb rd, off(rs1) (a sign-extended byte).
func (p *Program) Lb(rd, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("lb", err)
	}

	return p.instrLine("lb", p.b.Lb(rd, rs1, v), nil, pos)
}

// Sb - sb rs2, off(rs1) (the low byte of rs2).
func (p *Program) Sb(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("sb", err)
	}

	return p.instrLine("sb", p.b.Sb(rs2, rs1, v), nil, pos)
}

// Sw - sw rs2, off(rs1) (the low word of rs2).
func (p *Program) Sw(rs2, rs1 arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Off(off)
	if err != nil {
		return p.fail("sw", err)
	}

	return p.instrLine("sw", p.b.Sw(rs2, rs1, v), nil, pos)
}

// Ecall - ecall (the syscall trap).
func (p *Program) Ecall() *Program {
	pos := p.pos()
	i, err := arch.BuildInstr("ecall", nil)
	return p.instrLine("ecall", i, err, pos)
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
