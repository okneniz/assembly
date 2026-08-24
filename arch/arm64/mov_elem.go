package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MovElem - mov.d vd, vn[idx] (imm5[4] - the d-element index).
type MovElem struct {
	base

	rd, rn string
	idx    uint32
	enc    uint32
}

const movElemEnc uint32 = 0x4E003C00

func decodeMovElem(w uint32, addr uint64) Instr {
	return MovElem{
		base: newBase(addr, w),
		rd:   regNameX(w & 0x1f),
		rn:   vReg(w >> 5 & 0x1f),
		idx:  w >> 16 >> 4 & 1,
		enc:  movElemEnc,
	}
}

func (i MovElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mov.d %s, %s[%d]", i.rd, i.rn, i.idx)
}

func (i MovElem) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, err := armRegNum(i.rd)
	if err != nil {
		return 0, fmt.Errorf("mov.d: %w", err)
	}

	rn, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("mov.d: %w", err)
	}

	return writeWord(w, movElemEnc|rd|rn<<5|i.idx<<20)
}

func (i MovElem) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"mov",
		i.ObjDump(disasm.DefaultViewCtx()),
		"ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
