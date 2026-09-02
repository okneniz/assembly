package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrinsD - bstrins.d rd, rj, msb, lsb (DJUk6Um6): insert the rj[msb:lsb] field into rd[msb:lsb].
type BstrinsD struct {
	base

	rd, rj uint8
	msb    imm
	lsb    imm
}

// BstrinsD - bstrins.d rd, rj, msb, lsb.
func (Builder) BstrinsD(rd, rj Reg, msb, lsb UImm6) Instr {
	return BstrinsD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		msb: immNum(msb.Val()),
		lsb: immNum(lsb.Val()),
	}
}

func decodeBstrinsD(w uint32, addr uint64) Instr {
	return BstrinsD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		msb:  immNum(int64(uField(w, 16, 6))),
		lsb:  immNum(int64(uField(w, 10, 6))),
	}
}

func (i BstrinsD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrins.d %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrinsD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bstrins.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 6) | scatterU(i.lsb.val, 10, 6)

	return writeWord(w, word)
}

func (i BstrinsD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bstrins.d",
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
