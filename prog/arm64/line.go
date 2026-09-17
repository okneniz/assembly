package arm64

import (
	arch "github.com/okneniz/assembly/arch/arm64"
	"github.com/okneniz/assembly/prog"
)

// line - one source line. Exactly one of: a computed instruction (4
// bytes), a label-directed branch (4 bytes, the ctor receives the resolved
// absolute target and its own address), an La pair (8 bytes), a label
// definition, raw data, or a bss reserve. Every line belongs to a stream:
// the text stream (instructions, and read-only data before Data() switches
// away from it) or the data stream (Data() onwards, until Text()).
type line struct {
	instr  arch.Instr
	branch func(target, pc uint64) (arch.Instr, error)
	la     func(target, pc uint64) ([]arch.Instr, error)
	target string
	label  string
	data   []byte
	nobits int // a bss reserve: memory in the data stream, no file bytes
	stream int // 0: text, 1: data
	source string
	pos    prog.Pos
}

// size - the encoded size of the line (for the layout pass).
func (l line) size() int {
	switch {
	case l.label != "":
		return 0
	case l.nobits > 0:
		return l.nobits
	case l.data != nil:
		return len(l.data)
	case l.la != nil:
		return 8
	default:
		return 4
	}
}

func newLabelLine(name string, stream int, pos prog.Pos) line {
	return line{
		label:  name,
		stream: stream,
		source: name + ":",
		pos:    pos,
	}
}

func newDataLine(data []byte, src string, stream int, pos prog.Pos) line {
	return line{
		data:   data,
		stream: stream,
		source: src,
		pos:    pos,
	}
}

// newBssLine - a zero-fill reserve of the data stream.
func newBssLine(reserve int, pos prog.Pos) line {
	return line{
		nobits: reserve,
		stream: 1,
		source: ".bss",
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
