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
