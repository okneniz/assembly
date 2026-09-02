package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BytepickW - bytepick.w rd, rj, rk, sel (DJKUa2): pick 4 bytes out of the
// {rj, rk} concatenation, byte-indexed by sel.
type BytepickW struct {
	base

	rd, rj, rk uint8
	sel        imm
}

// BytepickW - bytepick.w rd, rj, rk, sel.
func (Builder) BytepickW(rd, rj, rk Reg, sel UImm2) Instr {
	return BytepickW{
		rd:  rd.Num(),
		rj:  rj.Num(),
		rk:  rk.Num(),
		sel: immNum(sel.Val()),
	}
}

func decodeBytepickW(w uint32, addr uint64) Instr {
	return BytepickW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		sel:  immNum(int64(uField(w, 15, 2))),
	}
}

func (i BytepickW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bytepick.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		laRegName(i.rk),
		i.sel.text(),
	)
}

func (i BytepickW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bytepick.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10 |
		scatterU(i.sel.val, 15, 2)

	return writeWord(w, word)
}

func (i BytepickW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bytepick.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{
			"rd":  laRegName(i.rd),
			"rj":  laRegName(i.rj),
			"rk":  laRegName(i.rk),
			"sel": i.sel.val,
		},
	)
}
