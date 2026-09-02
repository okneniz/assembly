package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LslReg — lsl rd, rn, rm.
type LslReg struct {
	base

	rd, rn, rm string
}

const LslRegX uint32 = 0x9A002000

// LslReg — lsl rd, rn, rm. Register 31 reads as zr (SP/WSP are
// not allowed — use XZR/WZR); the width is shared by all three registers.
func (Builder) LslReg(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "LslReg", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "LslReg", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "LslReg", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("LslReg", rd, rn, rm); err != nil {
		return nil, err
	}

	return LslReg{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

func decodeLslReg(w uint32, addr uint64) Instr {
	return LslReg{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i LslReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lsl %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i LslReg) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, LslRegX, 0x1A002000)
	if err != nil {
		return 0, fmt.Errorf("lsl: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("lsl: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i LslReg) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"lsl",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
