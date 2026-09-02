package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AsrReg — asr rd, rn, rm.
type AsrReg struct {
	base

	rd, rn, rm string
}

const AsrRegX uint32 = 0x9A002800

// AsrReg — asr rd, rn, rm. Register 31 reads as zr (SP/WSP are
// not allowed — use XZR/WZR); the width is shared by all three registers.
func (Builder) AsrReg(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "AsrReg", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AsrReg", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "AsrReg", "rm", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("AsrReg", rd, rn, rm); err != nil {
		return nil, err
	}

	return AsrReg{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

func decodeAsrReg(w uint32, addr uint64) Instr {
	return AsrReg{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i AsrReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asr %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i AsrReg) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, AsrRegX, 0x1A002800)
	if err != nil {
		return 0, fmt.Errorf("asr: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("asr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i AsrReg) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"asr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
