package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrpickW - bstrpick.w rd, rj, msb, lsb (DJUk5Um5): extract the rj[msb:lsb] field into rd, zero-extended.
type BstrpickW struct {
	base

	rd, rj uint8
	msb    imm
	lsb    imm
}

// BstrpickW - bstrpick.w rd, rj, msb, lsb.
func (Builder) BstrpickW(rd, rj Reg, msb, lsb UImm5) Instr {
	return BstrpickW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

func decodeBstrpickW(w uint32) Instr {
	return BstrpickW{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 5))),
		lsb:  immNum(int64(uField(w, 10, 5))),
	}
}

func (i BstrpickW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrpick.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrpickW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bstrpick.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 5) | scatterU(i.lsb.val, 10, 5)

	return writeWord(w, word)
}
