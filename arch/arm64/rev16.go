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

const Rev16X uint32 = 0xDAC00400

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

	return Rev16{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

func decodeRev16(w uint32, addr uint64) Instr {
	return Rev16{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
	}
}

func (i Rev16) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev16 %s, %s", i.rd, i.rn)
}

func (i Rev16) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i Rev16) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"rev16",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
