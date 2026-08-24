package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Uaddlv - uaddlv.Arr hN/sN/dN, vn (scalar dest by size).
type Uaddlv struct {
	base

	rd, rn  string
	q, size uint32
}

const uaddlvEnc uint32 = 0x2E303800

func decodeUaddlv(w uint32, addr uint64) Instr {
	rd := vReg(w & 0x1f)
	size := w >> 22 & 3
	switch size {
	case 0:
		rd = fmt.Sprintf("h%d", regIndex(rd))
	case 1:
		rd = fmt.Sprintf("s%d", regIndex(rd))
	}

	return Uaddlv{
		base: newBase(addr, w),
		rd:   rd,
		rn:   vReg(w >> 5 & 0x1f),
		q:    w >> 30 & 1,
		size: size,
	}
}

func (i Uaddlv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("uaddlv.%s %s, %s", decodeArrangement(i.q, i.size), i.rd, i.rn)
}

func (i Uaddlv) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, err := regNums2(fmt.Sprintf("v%d", regIndex(i.rd)), i.rn)
	if err != nil {
		return 0, fmt.Errorf("uaddlv: %w", err)
	}

	return writeWord(w, uaddlvEnc|rd|rn<<5)
}

func (i Uaddlv) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"uaddlv",
		i.ObjDump(disasm.DefaultViewCtx()),
		"ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
