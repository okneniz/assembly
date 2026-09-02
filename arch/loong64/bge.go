package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bge - bge rj, rd, offs (2RI16, the manual order swaps the registers):
// if the signed rj >= rd, jump to pc + offs (offs is word-scaled). The
// decoded form stores the absolute target.
type Bge struct {
	base

	rd, rj uint8
	off    imm
}

// Bge - bge rj, rd, offs (the pc-relative byte offset).
func (Builder) Bge(rj, rd Reg, off int64) Instr {
	return Bge{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeBge(w uint32, addr uint64) Instr {
	return Bge{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func (i Bge) ObjDump(_ disasm.ViewCtx) string {
	if i.rj == 0 {
		return fmt.Sprintf("blez %s, %s", laRegName(i.rd), i.off.text())
	}

	if i.rd == 0 {
		return fmt.Sprintf("bgez %s, %s", laRegName(i.rj), i.off.text())
	}

	return fmt.Sprintf("bge %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}

func (i Bge) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 16, "bge offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bge"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}

func (i Bge) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bge",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "rd": laRegName(i.rd), "off": i.off.val},
	)
}
