package arm64

import (
	"bytes"
	"fmt"

	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

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
// result: the code, the symbol table, the line map, and the assembly
// errors (undefined labels and entries, encode failures).
func (b *Binary) Assemble(base uint64) *prog.Result {
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

	// pass 2: encode; the line map grows with the pc.
	var lines []prog.LineEntry
	var code bytes.Buffer
	code.Grow(off)
	pc := base
	for _, l := range b.lines {
		if l.label == "" {
			lines = append(lines, prog.NewLineEntry(pc, l.size(), l.pos))
		}

		switch {
		case l.label != "":
			// nothing to encode
		case l.data != nil:
			_, _ = code.Write(l.data) // a bytes.Buffer write never fails
			pc += uint64(len(l.data))
		default:
			i, err := b.materialize(l, syms, pc)
			if err == nil {
				_, err = i.Encode(&code)
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

	return prog.NewResult(code.Bytes(), syms, lines, errs)
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
