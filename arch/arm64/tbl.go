package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tbl — tbl.16b vd, { vn }, vm.
type Tbl struct {
	base

	rd, rn, rm string
}

const tblEnc uint32 = 0x0E000000

func decodeTbl(w uint32, addr uint64) Instr {
	return Tbl{
		base: newBase(addr, w),
		rd:   vReg(w & 0x1f),
		rn:   vReg(w >> 5 & 0x1f),
		rm:   vReg(w >> 16 & 0x1f),
	}
}

func (i Tbl) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("tbl.16b %s, { %s }, %s", i.rd, i.rn, i.rm)
}

func (i Tbl) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("tbl: %w", err)
	}

	return writeWord(w, tblEnc|rd|rn<<5|rm<<16)
}

func (i Tbl) MarshalJSON() ([]byte, error) {
	return i.marshal("tbl", i.ObjDump(disasm.DefaultViewCtx()), "ASIMD",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm})
}
