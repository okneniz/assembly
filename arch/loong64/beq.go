package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Beq - beq rj, rd, offs (2RI16, the manual order swaps the registers):
// if rj == rd, jump to pc + offs (offs is word-scaled). The decoded form
// stores the absolute target.
type Beq struct {
	base

	rd, rj uint8
	off    imm
}

// NewBeq - beq rj, rd, offs (the pc-relative byte offset).
func NewBeq(rj, rd Reg, off int64) Instr {
	return Beq{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeBeq(w uint32, addr uint64) Instr {
	return Beq{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func (i Beq) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("beq %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}

func (i Beq) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 16, "beq offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["beq"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}

func (i Beq) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"beq",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "rd": laRegName(i.rd), "off": i.off.val},
	)
}
