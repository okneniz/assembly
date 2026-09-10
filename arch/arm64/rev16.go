package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rev16 — rev16 rd, rn.
type Rev16 struct {
	base

	rd, rn string
}

// newRev16 - the Rev16 constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newRev16(b base, rd Reg, rn Reg) (Rev16, error) {
	err := requireClass(
		rd,
		"Rev16",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Rev16{}, err
	}

	err = requireClass(
		rn,
		"Rev16",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Rev16{}, err
	}

	err = requireWidth(
		"Rev16",
		rd,
		rn,
	)

	if err != nil {
		return Rev16{}, err
	}

	return Rev16{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const Rev16X uint32 = 0xDAC00400

func (i Rev16) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev16 %s, %s", i.rd, i.rn)
}

func (i Rev16) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, Rev16X, 0x5AC00400)
	if err != nil {
		return 0, fmt.Errorf("rev16: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rev16: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Rev16(rd, rn Reg) (Instr, error) {
	return newRev16(base{}, rd, rn)
}

func decodeRev16(w uint32) (Instr, error) {
	in, err := newRev16(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
