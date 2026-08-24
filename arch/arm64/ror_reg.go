package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RorReg — ror rd, rn, rm.
type RorReg struct {
	base

	rd, rn, rm string
}

const RorRegX uint32 = 0x9A002C00

func decodeRorReg(w uint32, addr uint64) Instr {
	return RorReg{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i RorReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ror %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i RorReg) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, RorRegX, 0)
	if err != nil {
		return 0, fmt.Errorf("ror: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ror: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i RorReg) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ror",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
