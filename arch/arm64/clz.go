package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Clz — clz rd, rn.
type Clz struct {
	base

	rd, rn string
}

const ClzX uint32 = 0xDAC01000

// Clz — clz rd, rn. Register 31 reads as zr (SP/WSP are not
// allowed — use XZR/WZR); the width is shared by both registers.
func (Builder) Clz(rd, rn Reg) (Instr, error) {
	if err := requireClass(rd, "Clz", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Clz", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("Clz", rd, rn); err != nil {
		return nil, err
	}

	return Clz{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

func decodeClz(w uint32, addr uint64) Instr {
	return Clz{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
	}
}

func (i Clz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("clz %s, %s", i.rd, i.rn)
}

func (i Clz) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, ClzX, 0x5AC01000)
	if err != nil {
		return 0, fmt.Errorf("clz: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("clz: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (i Clz) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"clz",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
