package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BytepickD - bytepick.d rd, rj, rk, sel (DJKUa3): pick 8 bytes out of the
// {rj, rk} concatenation, byte-indexed by sel.
type BytepickD struct {
	base

	rd, rj, rk uint8
	sel        imm
}

// NewBytepickD - bytepick.d rd, rj, rk, sel.
func NewBytepickD(rd, rj, rk Reg, sel UImm3) Instr {
	return BytepickD{
		rd:  rd.Num(),
		rj:  rj.Num(),
		rk:  rk.Num(),
		sel: immNum(sel.Val()),
	}
}

func decodeBytepickD(w uint32, addr uint64) Instr {
	return BytepickD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
		sel:  immNum(int64(uField(w, 15, 3))),
	}
}

func (i BytepickD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bytepick.d %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		laRegName(i.rk),
		i.sel.text(),
	)
}

func (i BytepickD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bytepick.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10 |
		scatterU(i.sel.val, 15, 3)

	return writeWord(w, word)
}

func (i BytepickD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bytepick.d",
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
