package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ccmp — ccmp rn, rm, #imm, cond.
type Ccmp struct {
	base

	rn, rm string
	immVal uint32
	cond   string
}

const ccmpX uint32 = 0xFA400000

func decodeCcmp(w uint32, addr uint64) Instr {
	return Ccmp{
		base:   newBase(addr, w),
		rn:     armRegName(w>>5&0x1f, true),
		rm:     armRegName(w>>16&0x1f, true),
		immVal: w & 0xf,
		cond:   condName(w >> 12 & 0xf),
	}
}

func (i Ccmp) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ccmp %s, %s, #0x%x, %s", i.rn, i.rm, i.immVal, i.cond)
}

func (i Ccmp) Encode(w io.Writer, pc uint64) (int64, error) {
	rn, rm, err := regNums2(i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ccmp: %w", err)
	}

	c, err := condNum(i.cond)
	if err != nil {
		return 0, fmt.Errorf("ccmp: %w", err)
	}

	if i.immVal > 0xf {
		return 0, fmt.Errorf("ccmp: imm %#x out of range", i.immVal)
	}

	return writeWord(w, ccmpX|i.immVal|rn<<5|c<<12|rm<<16)
}

func (i Ccmp) MarshalJSON() ([]byte, error) {
	return i.marshal("ccmp", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rn": i.rn, "Rm": i.rm, "imm": i.immVal, "cond": i.cond})
}
