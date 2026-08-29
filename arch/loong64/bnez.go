package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bnez - bnez rj, offs (1RI21, the immediate split d5k16): if rj != 0,
// jump to pc + offs (word-scaled). The decoded form stores the absolute
// target.
type Bnez struct {
	base

	rj  uint8
	off imm
}

// NewBnez - bnez rj, offs (the pc-relative byte offset).
func NewBnez(rj Reg, off int64) Instr {
	return Bnez{
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeBnez(w uint32, addr uint64) Instr {
	return Bnez{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(d5k16Imm(w) << 2),
	}
}

func (i Bnez) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bnez %s, %s", laRegName(i.rj), i.off.text())
}

func (i Bnez) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 21, "bnez offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bnez"][0] | uint32(i.rj)<<5 | scatterD5k16(off)

	return writeWord(w, word)
}

func (i Bnez) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bnez",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "off": i.off.val},
	)
}
