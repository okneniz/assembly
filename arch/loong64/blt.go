package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Blt - blt rj, rd, offs (2RI16, the manual order swaps the registers):
// if the signed rj < rd, jump to pc + offs (offs is word-scaled). The
// decoded form stores the absolute target.
type Blt struct {
	base

	rd, rj uint8
	off    imm
}

// NewBlt - blt rj, rd, offs (the pc-relative byte offset).
func NewBlt(rj, rd Reg, off int64) Instr {
	return Blt{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeBlt(w uint32, addr uint64) Instr {
	return Blt{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func (i Blt) ObjDump(_ disasm.ViewCtx) string {
	if i.rd == 0 {
		return fmt.Sprintf("bltz %s, %s", laRegName(i.rj), i.off.text())
	}

	if i.rj == 0 {
		return fmt.Sprintf("bgtz %s, %s", laRegName(i.rd), i.off.text())
	}

	return fmt.Sprintf("blt %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}

func (i Blt) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 16, "blt offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["blt"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}

func (i Blt) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"blt",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "rd": laRegName(i.rd), "off": i.off.val},
	)
}
