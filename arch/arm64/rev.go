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

// newRev - the Rev constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newRev(b base, rd string, rn string) Rev {
	return Rev{
		base: b,
		rd:   rd,
		rn:   rn,
	}
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

// Rev — rev rd, rn. Only the 64-bit form: the package's Encode
// does not emit the 32-bit rev word (that encoding decodes as rev w,
// see the schemas).
func (Builder) Rev(rd, rn Reg) (Instr, error) {
	if err := requireClass(
		rd,
		"Rev",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Rev",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return newRev(base{}, rd.name(), rn.name()), nil
}

func decodeRev(w uint32) Instr {
	return newRev(newBase(w), armRegName(w&0x1f, w>>31&1 == 1), armRegName(w>>5&0x1f, w>>31&1 == 1))
}
