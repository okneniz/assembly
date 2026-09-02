// Package loong64 - programs written directly in Go: a chain Program
// over the arch builders, labels resolved at assembly time. The
// counterpart of an .s source file; immediate operands are raw ints,
// validation errors are deferred to Assemble.
package loong64

import (
	"fmt"

	arch "github.com/okneniz/assembly/arch/loong64"
)

// Program - a program being built: a sequence of lines (instructions,
// label-directed branches, la pairs, labels, data). Chain methods append
// one line each and return the program; nothing is encoded until
// Assemble.
type Program struct {
	lines []line
	entry string
	errs  []error
	b     arch.Builder
}

// line - one source line. Exactly one of: a computed instruction (4
// bytes), a label-directed branch (4 bytes, the ctor receives the
// resolved absolute target and its own address), an la pair (8 bytes),
// a label definition, or raw data.
type line struct {
	instr  arch.Instr
	branch func(target, pc uint64) (arch.Instr, error)
	la     func(target, pc uint64) ([]arch.Instr, error)
	target string
	label  string
	data   []byte
	source string
}

// size - the encoded size of the line (for the layout pass).
func (l line) size() int {
	switch {
	case l.label != "":
		return 0
	case l.data != nil:
		return len(l.data)
	case l.la != nil:
		return 8
	default:
		return 4
	}
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

// B - unconditional branch to a label (a pc-relative byte offset).
func (p *Program) B(label string) *Program {
	return p.branchLine("b", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.B(int64(t) - int64(pc)), nil
	})
}

// Beq - branch to a label when the registers are equal.
func (p *Program) Beq(rj, rd arch.Reg, label string) *Program {
	return p.branchLine("beq", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Beq(rj, rd, int64(t)-int64(pc)), nil
	})
}

// Bnez - branch to a label when the register is not zero.
func (p *Program) Bnez(rj arch.Reg, label string) *Program {
	return p.branchLine("bnez", label, func(t, pc uint64) (arch.Instr, error) {
		return p.b.Bnez(rj, int64(t)-int64(pc)), nil
	})
}

// La - load the address of a label into the register: the
// pcalau12i+addi.d pair (a fixed 8 bytes; the split is computed against
// the page-aligned pc, exactly as the text-path pseudo).
func (p *Program) La(rd arch.Reg, label string) *Program {
	p.lines = append(p.lines, line{
		la: func(t, pc uint64) ([]arch.Instr, error) {
			return laPair(p.b, rd, int64(t), int64(pc))
		},
		target: label,
		source: "la",
	})
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
	v, err := p.b.Imm20(imm)
	if err != nil {
		return p.fail("lu12i.w", err)
	}

	return p.instrLine("lu12i.w", p.b.Lu12iW(rd, v), nil)
}

// Ori - ori rd, rj, imm (zero-extended 12-bit).
func (p *Program) Ori(rd, rj arch.Reg, imm int64) *Program {
	v, err := p.b.UImm12(imm)
	if err != nil {
		return p.fail("ori", err)
	}

	return p.instrLine("ori", p.b.Ori(rd, rj, v), nil)
}

// LdBu - ld.bu rd, rj, off (a byte, zero-extended).
func (p *Program) LdBu(rd, rj arch.Reg, off int64) *Program {
	v, err := p.b.Imm12(off)
	if err != nil {
		return p.fail("ld.bu", err)
	}

	return p.instrLine("ld.bu", p.b.LdBu(rd, rj, v), nil)
}

// StB - st.b rd, rj, off (the low byte of rd).
func (p *Program) StB(rd, rj arch.Reg, off int64) *Program {
	v, err := p.b.Imm12(off)
	if err != nil {
		return p.fail("st.b", err)
	}

	return p.instrLine("st.b", p.b.StB(rd, rj, v), nil)
}

// AddiW - addi.w rd, rj, imm.
func (p *Program) AddiW(rd, rj arch.Reg, imm int64) *Program {
	v, err := p.b.Imm12(imm)
	if err != nil {
		return p.fail("addi.w", err)
	}

	return p.instrLine("addi.w", p.b.AddiW(rd, rj, v), nil)
}

// --- building and assembling ----------------------------------------------------

// Binary - a built program: the line sequence and the entry label, ready
// to be assembled (or transformed) at any base.
type Binary struct {
	Entry string
	lines []line
}

// Build - materialize the program; deferred construction errors are
// returned alongside.
func (p *Program) Build() (*Binary, []error) {
	return &Binary{Entry: p.entry, lines: p.lines}, p.errs
}

// Assemble - encode the program at base: labels become absolute
// addresses, label-directed lines receive their targets. Returns the
// code, the symbol table, and the assembly errors.
func (b *Binary) Assemble(base uint64) ([]byte, map[string]uint64, []error) {
	var errs []error

	// pass 1: layout.
	off := 0
	syms := make(map[string]uint64)
	for _, l := range b.lines {
		if l.label != "" {
			syms[l.label] = base + uint64(off)
		}

		off += l.size()
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
		case l.la != nil:
			target, ok := syms[l.target]
			if !ok {
				errs = append(errs, fmt.Errorf("%s: undefined label %q", l.source, l.target))
				pc += 8
				continue
			}

			pair, err := l.la(target, pc)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", l.source, err))
				pc += 8
				continue
			}

			for _, i := range pair {
				if _, err := i.Encode(byteWriter{&code}, pc); err != nil {
					errs = append(errs, fmt.Errorf("%s: %w", l.source, err))
				}

				pc += 4
			}
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

// materialize - the instruction of a line at its own address.
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

func (p *Program) instrLine(src string, i arch.Instr, _ error) *Program {
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
