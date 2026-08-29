package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrinsW - bstrins.w rd, rj, msb, lsb (DJUk5Um5): insert the rj[msb:lsb]
// field into rd[msb:lsb] (msb at [20:16], lsb at [14:10]).
type BstrinsW struct {
	base

	rd, rj uint8
	msb    imm
	lsb    imm
}

// NewBstrinsW - bstrins.w rd, rj, msb, lsb.
func NewBstrinsW(rd, rj Reg, msb, lsb UImm5) Instr {
	return BstrinsW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

func decodeBstrinsW(w uint32, addr uint64) Instr {
	return BstrinsW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 5))),
		lsb:  immNum(int64(uField(w, 10, 5))),
	}
}

func (i BstrinsW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrins.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrinsW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bstrins.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 5) | scatterU(i.lsb.val, 10, 5)

	return writeWord(w, word)
}

func (i BstrinsW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bstrins.w",
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
