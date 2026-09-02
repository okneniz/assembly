package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AlslD - alsl.d rd, rj, rk, shift (DJKUa2): rd = rj + (rk << shift), the
// 64-bit sum; the 1..4 shift is encoded as ui2 = shift - 1.
type AlslD struct {
	base

	rd, rj, rk uint8
	shift      imm
}

// AlslD - alsl.d rd, rj, rk, shift.
func (Builder) AlslD(rd, rj, rk Reg, shift Shift3) Instr {
	return AlslD{
		rd:    rd.Num(),
		rj:    rj.Num(),
		rk:    rk.Num(),
		shift: immNum(shift.Val()),
	}
}

func decodeAlslD(w uint32, addr uint64) Instr {
	return AlslD{
		base:  newBase(addr, w),
		rd:    uint8(w & 0x1f),
		rj:    uint8(w >> 5 & 0x1f),
		rk:    uint8(w >> 10 & 0x1f),
		shift: immNum(int64(uField(w, 15, 2)) + 1),
	}
}

func (i AlslD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"alsl.d %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		laRegName(i.rk),
		i.shift.text(),
	)
}

func (i AlslD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["alsl.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10 |
		scatterU(i.shift.val-1, 15, 2)

	return writeWord(w, word)
}

func (i AlslD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"alsl.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{
			"rd":    laRegName(i.rd),
			"rj":    laRegName(i.rj),
			"rk":    laRegName(i.rk),
			"shift": i.shift.val,
		},
	)
}
