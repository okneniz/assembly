package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmminDbWu - ammin_db.wu rd, rk, rj (DJK, @orig_fmt=DKJ - the second assembly
// operand packs the K field [14:10], the third the J field [9:5]):
// rd = AMO(MEM[rj], rk): unsigned fetch-min with a built-in barrier, old value to rd.
type AmminDbWu struct {
	base

	rd, rk, rj uint8
}

// NewAmminDbWu - ammin_db.wu rd, rk, rj (the operand order the assembler writes).
func NewAmminDbWu(rd, rk, rj Reg) Instr {
	return AmminDbWu{
		rd: rd.Num(),
		rk: rk.Num(),
		rj: rj.Num(),
	}
}

func decodeAmminDbWu(w uint32, addr uint64) Instr {
	return AmminDbWu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i AmminDbWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ammin_db.wu %s, %s, %s", laRegName(i.rd), laRegName(i.rk), laRegName(i.rj))
}

func (i AmminDbWu) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ammin_db.wu"][0] |
		uint32(i.rd) | uint32(i.rk)<<10 | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i AmminDbWu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ammin_db.wu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rk": laRegName(i.rk), "rj": laRegName(i.rj)},
	)
}
