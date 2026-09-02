// Package arm64 - programs written directly in Go: a chain Program over
// the arch builders, labels resolved at assembly time. The counterpart of
// an .s source file: every chain method is a source line, a macro is an
// ordinary Go function returning *Program. Immediate operands are raw
// ints; validation errors are deferred to Assemble.
package arm64

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
)

// Program - a program being built: a sequence of lines (instructions,
// label-directed branches, labels, data). Chain methods append one line
// each and return the program; nothing is encoded until Assemble.
type Program struct {
	lines []line
	entry string
	errs  []error
	b     arch.Builder
}

// line - one source line: exactly one of a computed instruction, a
// label-directed instruction (the ctor receives the resolved absolute
// target and its own address), a label definition, or raw data.
type line struct {
	instr  arch.Instr
	branch func(target, pc uint64) (arch.Instr, error)
	target string
	label  string
	data   []byte
	source string
}

// New - an empty program.
func New() *Program {
	return &Program{}
}

// Label - define a label at the current position.
func (p *Program) Label(name string) *Program {
	p.lines = append(p.lines, line{label: name, source: name + ":"})
	return p
}

// Entry - the label the program starts at.
func (p *Program) Entry(name string) *Program {
	p.entry = name
	return p
}

// Ascii - string data appended verbatim (no terminating zero).
func (p *Program) Ascii(s string) *Program {
	p.lines = append(p.lines, line{data: []byte(s), source: ".ascii"})
	return p
}

// Bytes - raw data bytes.
func (p *Program) Bytes(b ...byte) *Program {
	p.lines = append(p.lines, line{data: b, source: ".byte"})
	return p
}

// Instr - append an already computed instruction (the escape hatch when
// no chain twin exists yet).
func (p *Program) Instr(i arch.Instr, src string) *Program {
	p.lines = append(p.lines, line{instr: i, source: src})
	return p
}

// --- label-directed lines -------------------------------------------------------

// B - branch to a label.
func (p *Program) B(label string) *Program {
	return p.branchLine("b", label, func(t, _ uint64) (arch.Instr, error) {
		return p.b.B(int64(t)), nil
	})
}

// Bl - branch with link to a label.
func (p *Program) Bl(label string) *Program {
	return p.branchLine("bl", label, func(t, _ uint64) (arch.Instr, error) {
		return p.b.Bl(int64(t)), nil
	})
}

// Bcond - conditional branch to a label.
func (p *Program) Bcond(cond, label string) *Program {
	return p.branchLine("b."+cond, label, func(t, _ uint64) (arch.Instr, error) {
		return p.b.Bcond(cond, int64(t))
	})
}

// Cbz - branch to a label when the register is zero.
func (p *Program) Cbz(rt arch.Reg, label string) *Program {
	return p.branchLine("cbz", label, func(t, _ uint64) (arch.Instr, error) {
		return p.b.Cbz(rt, int64(t))
	})
}

// Cbnz - branch to a label when the register is not zero.
func (p *Program) Cbnz(rt arch.Reg, label string) *Program {
	return p.branchLine("cbnz", label, func(t, _ uint64) (arch.Instr, error) {
		return p.b.Cbnz(rt, int64(t))
	})
}

// Adr - load the address of a label into the register (pc-relative).
func (p *Program) Adr(rd arch.Reg, label string) *Program {
	return p.branchLine("adr", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Adr(rd, int64(t)-int64(pc))
	})
}

// --- instruction chain twins ----------------------------------------------------

// Movz - movz rd, #imm[, lsl #hw*16].
func (p *Program) Movz(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movz", err)
	}

	i, err := p.b.Movz(rd, v, hw)
	return p.instrLine("movz", i, err)
}

// Movk - movk rd, #imm, lsl #hw*16.
func (p *Program) Movk(rd arch.Reg, imm int64, hw arch.Hw) *Program {
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("movk", err)
	}

	i, err := p.b.Movk(rd, v, hw)
	return p.instrLine("movk", i, err)
}

// Mov - mov rd, #imm (the assembler's immediate form: movz; values wider
// than 16 bits come from Movz/Movk pairs - see the syscall constants).
func (p *Program) Mov(rd arch.Reg, imm int64) *Program {
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("mov", err)
	}

	i, err := p.b.Movz(rd, v, arch.Hw0)
	return p.instrLine("mov", i, err)
}

// Svc - svc #imm (the canonical Darwin trap immediate is 0x80).
func (p *Program) Svc(imm int64) *Program {
	v, err := p.b.Imm16(imm)
	if err != nil {
		return p.fail("svc", err)
	}

	return p.instrLine("svc", p.b.Svc(v), nil)
}

// --- building and assembling ----------------------------------------------------

// Binary - a built program: the line sequence and the entry label, ready
// to be assembled (or transformed) at any base.
type Binary struct {
	Entry string
	lines []line
}

// Build - materialize the program; deferred construction errors are
// returned alongside (an empty slice means the program is well-formed).
func (p *Program) Build() (*Binary, []error) {
	return &Binary{Entry: p.entry, lines: p.lines}, p.errs
}

// Assemble - encode the program at base: labels become absolute
// addresses, label-directed lines receive their targets. Returns the
// code, the symbol table, and the assembly errors (undefined labels and
// entries, encode failures).
func (b *Binary) Assemble(base uint64) ([]byte, map[string]uint64, []error) {
	var errs []error

	// pass 1: layout - fixed 4-byte instructions, data lines verbatim.
	off := 0
	syms := make(map[string]uint64)
	for _, l := range b.lines {
		switch {
		case l.label != "":
			syms[l.label] = base + uint64(off)
		case l.data != nil:
			off += len(l.data)
		default:
			off += 4
		}
	}

	// pass 2: encode.
	code := make([]byte, 0, off)
	pc := base
	for _, l := range b.lines {
		switch {
		case l.label != "":
			// nothing to encode
		case l.data != nil:
			code = append(code, l.data...)
			pc += uint64(len(l.data))
		default:
			i, err := b.materialize(l, syms, pc)
			if err == nil {
				_, err = i.Encode(byteWriter{&code}, pc)
			}

			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", l.source, err))
			}

			pc += 4
		}
	}

	if b.Entry != "" {
		if _, ok := syms[b.Entry]; !ok {
			errs = append(errs, fmt.Errorf("entry: undefined label %q", b.Entry))
		}
	}

	return code, syms, errs
}

// materialize - the instruction of a line at its own address: computed
// lines pass through, label-directed lines get their target resolved.
func (b *Binary) materialize(l line, syms map[string]uint64, pc uint64) (arch.Instr, error) {
	if l.branch == nil {
		return l.instr, nil
	}

	target, ok := syms[l.target]
	if !ok {
		return nil, fmt.Errorf("undefined label %q", l.target)
	}

	return l.branch(target, pc)
}

// --- internals ---------------------------------------------------------------

func (p *Program) instrLine(src string, i arch.Instr, err error) *Program {
	if err != nil {
		return p.fail(src, err)
	}

	p.lines = append(p.lines, line{instr: i, source: src})
	return p
}

func (p *Program) branchLine(
	src, label string,
	ctor func(target, pc uint64) (arch.Instr, error),
) *Program {
	p.lines = append(p.lines, line{branch: ctor, target: label, source: src})
	return p
}

func (p *Program) fail(src string, err error) *Program {
	p.errs = append(p.errs, fmt.Errorf("%s: %w", src, err))
	return p
}

// byteWriter - an io.Writer appending to a byte slice.
type byteWriter struct{ b *[]byte }

func (w byteWriter) Write(p []byte) (int, error) {
	*w.b = append(*w.b, p...)
	return len(p), nil
}
