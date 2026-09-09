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

// newRev32 - the Rev32 constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newRev32(b base, rd string, rn string) Rev32 {
	return Rev32{
		base: b,
		rd:   rd,
		rn:   rn,
	}
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

// Rev32 — rev32 rd, rn. Only the 64-bit form: the architecture has
// no 32-bit rev32 (the sf=0 slot of the encoding is the 32-bit rev).
func (Builder) Rev32(rd, rn Reg) (Instr, error) {
	if err := requireClass(
		rd,
		"Rev32",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Rev32",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return newRev32(base{}, rd.name(), rn.name()), nil
}

func decodeRev32(w uint32) Instr {
	return newRev32(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
	)
}
