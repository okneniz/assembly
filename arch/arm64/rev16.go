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

// newRev16 - the Rev16 constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newRev16(b base, rd string, rn string) Rev16 {
	return Rev16{
		base: b,
		rd:   rd,
		rn:   rn,
	}
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

// Rev16 — rev16 rd, rn. Register 31 reads as zr (SP/WSP are not
// allowed — use XZR/WZR); the width is shared by both registers.
func (Builder) Rev16(rd, rn Reg) (Instr, error) {
	if err := requireClass(rd, "Rev16", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Rev16", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("Rev16", rd, rn); err != nil {
		return nil, err
	}

	return newRev16(base{}, rd.name(), rn.name()), nil
}

func decodeRev16(w uint32) Instr {
	return newRev16(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
	)
}
