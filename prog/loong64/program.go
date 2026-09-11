// Package loong64 - programs written directly in Go: a chain Program
// over the arch builders, labels resolved at assembly time. The
// counterpart of an .s source file; immediate operands are raw ints,
// validation errors are deferred to Assemble. The Go source position of
// each line (the debugger's address ↔ source map) comes from the
// position resolver injected with WithPos - there is no built-in caller
// detection.
package loong64

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/loong64"
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

// B - unconditional branch to a label (a pc-relative byte offset).
func (p *Program) B(label string) *Program {
	pos := p.pos()
	return p.branchLine("b", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.B(int64(t) - int64(pc)), nil
	}, pos)
}

// Beq - branch to a label when the registers are equal.
func (p *Program) Beq(rj, rd arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("beq", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Beq(rj, rd, int64(t)-int64(pc)), nil
	}, pos)
}

// Bnez - branch to a label when the register is not zero.
func (p *Program) Bnez(rj arch.Reg, label string) *Program {
	pos := p.pos()
	return p.branchLine("bnez", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bnez(rj, int64(t)-int64(pc)), nil
	}, pos)
}

// La - load the address of a label into the register: the
// pcalau12i+addi.d pair (a fixed 8 bytes; the split is computed against
// the page-aligned pc, exactly as the text-path pseudo).
func (p *Program) La(rd arch.Reg, label string) *Program {
	pos := p.pos()
	p.lines = append(p.lines, newLaLine(label, func(t, pc uint64) ([]arch.Instr, error) {
		return laPair(p.b, rd, int64(t), int64(pc))
	}, pos))
	return p
}

// laPair - the evaluated la encoding: pcalau12i (hi) + addi.d (lo).
func laPair(b arch.Builder, rd arch.Reg, target, pc int64) ([]arch.Instr, error) {
	page := pc &^ 0xfff
	d := target - page
	lo := d & 0xfff
	if lo >= 0x800 {
		lo -= 0x1000
	}

	hi := (d - lo) >> 12
	hi20, err := b.Imm20(hi)
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	lo12, err := b.Imm12(lo)
	if err != nil {
		return nil, fmt.Errorf("la: %w", err)
	}

	return []arch.Instr{
		b.Pcalau12i(rd, hi20),
		b.AddiD(rd, rd, lo12),
	}, nil
}

// --- instruction chain twins ----------------------------------------------------

// Lu12iW - lu12i.w rd, imm (the high 20 bits of an address).
func (p *Program) Lu12iW(rd arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm20(imm)
	if err != nil {
		return p.fail("lu12i.w", err)
	}

	return p.instrLine("lu12i.w", p.b.Lu12iW(rd, v), nil, pos)
}

// Ori - ori rd, rj, imm (zero-extended 12-bit).
func (p *Program) Ori(rd, rj arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.UImm12(imm)
	if err != nil {
		return p.fail("ori", err)
	}

	return p.instrLine("ori", p.b.Ori(rd, rj, v), nil, pos)
}

// LdBu - ld.bu rd, rj, off (a byte, zero-extended).
func (p *Program) LdBu(rd, rj arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(off)
	if err != nil {
		return p.fail("ld.bu", err)
	}

	return p.instrLine("ld.bu", p.b.LdBu(rd, rj, v), nil, pos)
}

// StB - st.b rd, rj, off (the low byte of rd).
func (p *Program) StB(rd, rj arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(off)
	if err != nil {
		return p.fail("st.b", err)
	}

	return p.instrLine("st.b", p.b.StB(rd, rj, v), nil, pos)
}

// AddiW - addi.w rd, rj, imm.
func (p *Program) AddiW(rd, rj arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("addi.w", err)
	}

	return p.instrLine("addi.w", p.b.AddiW(rd, rj, v), nil, pos)
}

// --- internals ---------------------------------------------------------------

func (p *Program) instrLine(src string, i arch.Instr, _ error, pos prog.Pos) *Program {
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
