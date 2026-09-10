package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rev — rev rd, rn.
type Rev struct {
	base

	rd, rn string
}

// newRev - the Rev constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newRev(b base, rd Reg, rn Reg) (Rev, error) {
	err := requireClass(
		rd,
		"Rev",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Rev{}, err
	}

	err = requireClass(
		rn,
		"Rev",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Rev{}, err
	}

	return Rev{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const RevX uint32 = 0xDAC00C00

func (i Rev) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev %s, %s", i.rd, i.rn)
}

func (i Rev) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, RevX, 0)
	if err != nil {
		return 0, fmt.Errorf("rev: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rev: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Rev(rd, rn Reg) (Instr, error) {
	return newRev(base{}, rd, rn)
}

func decodeRev(w uint32) (Instr, error) {
	in, err := newRev(newBase(w), gprOf(w&0x1f, w>>31&1 == 1), gprOf(w>>5&0x1f, w>>31&1 == 1))
	if err != nil {
		return nil, err
	}

	return in, nil
}
