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
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

// Program - a program being built: a sequence of lines (instructions,
// label-directed branches, labels, data). Chain methods append one line
// each and return the program; nothing is encoded until Assemble.
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

// Mov - mov rd, #imm (the assembler's immediate form: movz; values wider
// than 16 bits come from Movz/Movk pairs - see the syscall constants).
func (p *Program) Mov(rd arch.Reg, imm int64) *Program {
	pos := p.pos()
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("mov", err)
	}

	i, err := p.b.Movz(rd, v, arch.Hw0)
	return p.instrLine("mov", i, err, pos)
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
