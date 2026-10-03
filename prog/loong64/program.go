// Package loong64 - programs written directly in Go: a chain Program
// over the arch builders, labels resolved at assembly time. The
// counterpart of an .s source file; immediate operands are raw ints,
// validation errors are deferred to Assemble. The Go source position of
// each line (the debugger's address ↔ source map) comes from the
// position resolver injected with WithPos - there is no built-in caller
// detection.
// The chain is a producer over the unit output: every method deposits
// its record there (instructions as ready records, label-directed macros
// as deferred ones), and Assemble is the unit's resolve phase adapted to
// the chain's result shape.
package loong64

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/unit"
)

// Program - a program being built: a sequence of lines (instructions,
// label-directed branches, la pairs, labels, data). Chain methods append
// one line each and return the program; nothing is encoded until
// Assemble.
type Program struct {
	u    *unit.Unit
	errs []error
	pos  func() unit.Pos
	b    arch.Builder
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

// AddiW - addi.w rd, rj, imm.
func (p *Program) AddiW(rd, rj arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("addi.w", err)
	}

	return p.instrLine(p.b.AddiW(rd, rj, v), nil, pos)
}

// Ascii - string data appended verbatim (no terminating zero).
func (p *Program) Ascii(s string) *Program {
	p.u.Ascii(p.pos(), s)
	return p
}

// laPair - the evaluated la encoding: pcalau12i (hi) + addi.d (lo).
func laPair(b arch.Builder, rd arch.Reg, target, pc int64) ([]unit.Resolved, error) {
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

	return []unit.Resolved{
		b.Pcalau12i(rd, hi20),
		b.AddiD(rd, rd, lo12),
	}, nil
}

// --- instruction chain twins ----------------------------------------------------

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

// Build - materialize the program; deferred construction errors are
// returned alongside.
func (p *Program) Build() (*Binary, []error) {
	return &Binary{u: p.u}, p.errs
}

// --- internals ---------------------------------------------------------------

// Bytes - raw data bytes.
func (p *Program) Bytes(b ...byte) *Program {
	p.u.Bytes(p.pos(), b...)
	return p
}

// --- label-directed lines -------------------------------------------------------

// Entry - the label the program starts at.
func (p *Program) Entry(name string) *Program {
	p.u.Entry(name)
	return p
}

// La - load the address of a label into the register: the
// pcalau12i+addi.d pair (a fixed 8 bytes; the split is computed against
// the page-aligned pc, exactly as the text-path pseudo).
func (p *Program) La(rd arch.Reg, label string) *Program {
	p.u.Sym(p.pos(), unit.NewPair("la", label, 8, func(t, pc uint64) ([]unit.Resolved, error) {
		return laPair(p.b, rd, int64(t), int64(pc))
	}))
	return p
}

// Label - define a label at the current position.
func (p *Program) Label(name string) *Program {
	p.u.Label(name)
	return p
}

// LdBu - ld.bu rd, rj, off (a byte, zero-extended).
func (p *Program) LdBu(rd, rj arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(off)
	if err != nil {
		return p.fail("ld.bu", err)
	}

	return p.instrLine(p.b.LdBu(rd, rj, v), nil, pos)
}

// Lu12iW - lu12i.w rd, imm (the high 20 bits of an address).
func (p *Program) Lu12iW(rd arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm20(imm)
	if err != nil {
		return p.fail("lu12i.w", err)
	}

	return p.instrLine(p.b.Lu12iW(rd, v), nil, pos)
}

// Ori - ori rd, rj, imm (zero-extended 12-bit).
func (p *Program) Ori(rd, rj arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.UImm12(imm)
	if err != nil {
		return p.fail("ori", err)
	}

	return p.instrLine(p.b.Ori(rd, rj, v), nil, pos)
}

// StB - st.b rd, rj, off (the low byte of rd).
func (p *Program) StB(rd, rj arch.Reg, off int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm12(off)
	if err != nil {
		return p.fail("st.b", err)
	}

	return p.instrLine(p.b.StB(rd, rj, v), nil, pos)
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

func (p *Program) branchLine(
	src, label string,
	ctor func(target, pc uint64) (arch.Instr, error),
	pos unit.Pos,
) *Program {
	p.u.Sym(pos, unit.NewBranch(src, label, 4, func(t, pc uint64) (unit.Resolved, error) {
		return ctor(t, pc)
	}))
	return p
}

func (p *Program) fail(src string, err error) *Program {
	p.errs = append(p.errs, fmt.Errorf("%s: %w", src, err))
	return p
}

func (p *Program) instrLine(i arch.Instr, _ error, pos unit.Pos) *Program {
	p.u.Instr(pos, i, nil)
	return p
}
