package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bgeu - bgeu rj, rd, offs (2RI16, the manual order swaps the registers):
// if the unsigned rj >= rd, jump to pc + offs (offs is word-scaled). The
// decoded form stores the absolute target.
type Bgeu struct {
	base

	rd, rj uint8
	off    imm
}

// NewBgeu - bgeu rj, rd, offs (the pc-relative byte offset).
func NewBgeu(rj, rd Reg, off int64) Instr {
	return Bgeu{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeBgeu(w uint32, addr uint64) Instr {
	return Bgeu{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func (i Bgeu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bgeu %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}

func (i Bgeu) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 16, "bgeu offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bgeu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}

func (i Bgeu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bgeu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "rd": laRegName(i.rd), "off": i.off.val},
	)
}
