package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cls — cls rd, rn.
type Cls struct {
	base

	rd, rn string
}

const ClsX uint32 = 0xDAC01400

// Cls — cls rd, rn. Register 31 reads as zr (SP/WSP are not
// allowed — use XZR/WZR); the width is shared by both registers.
func (Builder) Cls(rd, rn Reg) (Instr, error) {
	if err := requireClass(rd, "Cls", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Cls", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("Cls", rd, rn); err != nil {
		return nil, err
	}

	return Cls{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

func decodeCls(w uint32, addr uint64) Instr {
	return Cls{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
	}
}

func (i Cls) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cls %s, %s", i.rd, i.rn)
}

func (i Cls) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, ClsX, 0x5AC01400)
	if err != nil {
		return 0, fmt.Errorf("cls: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("cls: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (i Cls) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"cls",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
