package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmmaxDbDu - ammax_db.du rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): unsigned fetch-max with a built-in barrier, old value to rd.
type AmmaxDbDu struct {
	base

	rd, rk, rj uint8
}

// NewAmmaxDbDu - ammax_db.du rd, rk, rj (the operand order the assembler writes).
func NewAmmaxDbDu(rd, rk, rj Reg) Instr {
	return AmmaxDbDu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmmaxDbDu(w uint32, addr uint64) Instr {
	return AmmaxDbDu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmmaxDbDu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammax_db.du %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmmaxDbDu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammax_db.du"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmmaxDbDu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammax_db.du",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
