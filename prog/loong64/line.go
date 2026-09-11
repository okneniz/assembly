package loong64

import (
	arch "github.com/okneniz/assembly/arch/loong64"
	"github.com/okneniz/assembly/prog"
)

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
	pos    prog.Pos
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

func newLabelLine(name string, pos prog.Pos) line {
	return line{
		label:  name,
		source: name + ":",
		pos:    pos,
	}
}

func newDataLine(data []byte, src string, pos prog.Pos) line {
	return line{
		data:   data,
		source: src,
		pos:    pos,
	}
}

func newInstrLine(i arch.Instr, src string, pos prog.Pos) line {
	return line{
		instr:  i,
		source: src,
		pos:    pos,
	}
}

func newBranchLine(
	src, target string,
	ctor func(target, pc uint64) (arch.Instr, error),
	pos prog.Pos,
) line {
	return line{
		branch: ctor,
		target: target,
		source: src,
		pos:    pos,
	}
}

func newLaLine(
	target string,
	ctor func(target, pc uint64) ([]arch.Instr, error),
	pos prog.Pos,
) line {
	return line{
		la:     ctor,
		target: target,
		source: "la",
		pos:    pos,
	}
}
