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

// Assemble - encode the program at base as one flat stream: labels become
// absolute addresses, label-directed lines receive their targets. Returns
// the result: the code, the symbol table, the line map, and the assembly
// errors (undefined labels and entries, encode failures). A program that
// switched to the data stream needs AssembleLayout instead.
func (b *Binary) Assemble(base uint64) *prog.Result {
	var errs []error

	for _, l := range b.lines {
		if l.stream != 0 || l.nobits > 0 {
			errs = append(errs, fmt.Errorf(
				"%s: the data stream needs AssembleLayout", l.source,
			))
			return prog.NewResult(nil, nil, nil, errs)
		}
	}

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
		case l.la != nil:
			target, ok := syms[l.target]
			if !ok {
				errs = append(errs, fmt.Errorf(
					"%s: undefined label %q", l.source, l.target))
				pc += 8
				continue
			}

			pair, err := l.la(target, pc)
			if err == nil {
				for _, i := range pair {
					_, err = i.Encode(&code)
					if err != nil {
						break
					}
				}
			}

			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", l.source, err))
			}

			pc += 8
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

// labAt - a label position of the layout pass: its stream and the memory
// offset inside it.
type labAt struct {
	stream int
	off    uint64
}

// AssembleLayout - encode the program with its two streams placed by the
// policy: place sees the stream sizes (the data memory size includes the
// bss tail) and returns the base address of each stream. Labels resolve
// against those bases, so text may reference data statics across the page
// gap through the label-directed lines (La, Adrp, ...) without knowing
// the layout. file.MachoPlaceStreams is the Mach-O policy - the same
// numbers the writer places the bytes at.
func (b *Binary) AssembleLayout(place prog.Place) *prog.Result {
	var errs []error

	// pass 1: stream sizes and label positions.
	var mem, file [2]int
	labs := make(map[string]labAt)
	for _, l := range b.lines {
		if l.label != "" {
			labs[l.label] = labAt{stream: l.stream, off: uint64(mem[l.stream])}
		}

		mem[l.stream] += l.size()
		if l.nobits == 0 {
			file[l.stream] += l.size()
		}
	}

	textAddr, dataAddr := place(mem[0], file[1], mem[1])
	base := [2]uint64{textAddr, dataAddr}

	syms := make(map[string]uint64, len(labs))
	for name, at := range labs {
		syms[name] = base[at.stream] + at.off
	}

	// pass 2: encode per stream; the line map covers both.
	var lines []prog.LineEntry
	var code, data bytes.Buffer
	var cur [2]int

	for _, l := range b.lines {
		if l.label == "" {
			lines = append(lines, prog.NewLineEntry(
				base[l.stream]+uint64(cur[l.stream]), l.size(), l.pos))
		}

		pc := base[l.stream] + uint64(cur[l.stream])

		switch {
		case l.label != "":
			// nothing to encode
		case l.nobits > 0:
			cur[1] += l.nobits
		case l.data != nil:
			if l.stream == 0 {
				_, _ = code.Write(l.data)
			} else {
				_, _ = data.Write(l.data)
			}

			cur[l.stream] += len(l.data)
		case l.la != nil:
			target, ok := syms[l.target]
			if !ok {
				errs = append(errs, fmt.Errorf(
					"%s: undefined label %q", l.source, l.target))
				cur[0] += 8
				continue
			}

			pair, err := l.la(target, pc)
			if err == nil {
				for _, i := range pair {
					_, err = i.Encode(&code)
					if err != nil {
						break
					}
				}
			}

			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", l.source, err))
			}

			cur[0] += 8
		default:
			i, err := b.materialize(l, syms, pc)
			if err == nil {
				_, err = i.Encode(&code)
			}

			if err != nil {
				errs = append(errs, fmt.Errorf("%s: %w", l.source, err))
			}

			cur[0] += 4
		}
	}

	if b.Entry != "" {
		if _, ok := syms[b.Entry]; !ok {
			errs = append(errs, fmt.Errorf("entry: undefined label %q", b.Entry))
		}
	}

	res := prog.NewResult(code.Bytes(), syms, lines, errs)
	res.Data = data.Bytes()
	res.DataMem = mem[1]
	return res
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
