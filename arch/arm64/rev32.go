package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rev32 — rev32 rd, rn.
type Rev32 struct {
	base

	rd, rn string
}

// newRev32 - the Rev32 constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newRev32(b base, rd Reg, rn Reg) (Rev32, error) {
	err := requireClass(
		rd,
		"Rev32",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Rev32{}, err
	}

	err = requireClass(
		rn,
		"Rev32",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Rev32{}, err
	}

	return Rev32{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const Rev32X uint32 = 0xDAC00800

func (i Rev32) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev32 %s, %s", i.rd, i.rn)
}

func (i Rev32) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, Rev32X, 0x5AC00800)
	if err != nil {
		return 0, fmt.Errorf("rev32: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rev32: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Rev32(rd, rn Reg) (Instr, error) {
	return newRev32(base{}, rd, rn)
}

func decodeRev32(w uint32) (Instr, error) {
	in, err := newRev32(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
