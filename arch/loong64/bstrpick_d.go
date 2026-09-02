package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrpickD - bstrpick.d rd, rj, msb, lsb (DJUk6Um6): extract the rj[msb:lsb] field into rd, zero-extended.
type BstrpickD struct {
	base

	rd, rj uint8
	msb    imm
	lsb    imm
}

// BstrpickD - bstrpick.d rd, rj, msb, lsb.
func (Builder) BstrpickD(rd, rj Reg, msb, lsb UImm6) Instr {
	return BstrpickD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

func decodeBstrpickD(w uint32, addr uint64) Instr {
	return BstrpickD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 6))),
		lsb:  immNum(int64(uField(w, 10, 6))),
	}
}

func (i BstrpickD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrpick.d %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrpickD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bstrpick.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 6) | scatterU(i.lsb.val, 10, 6)

	return writeWord(w, word)
}

func (i BstrpickD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bstrpick.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{
			"rd":  laRegName(i.rd),
			"rj":  laRegName(i.rj),
			"msb": i.msb.val,
			"lsb": i.lsb.val,
		},
	)
}
