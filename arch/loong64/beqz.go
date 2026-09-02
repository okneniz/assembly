package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Beqz - beqz rj, offs (1RI21, the immediate split d5k16): if rj == 0,
// jump to pc + offs (word-scaled). The decoded form stores the absolute
// target.
type Beqz struct {
	base

	rj  uint8
	off imm
}

// Beqz - beqz rj, offs (the pc-relative byte offset).
func (Builder) Beqz(rj Reg, off int64) Instr {
	return Beqz{
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeBeqz(w uint32, addr uint64) Instr {
	return Beqz{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(d5k16Imm(w) << 2),
	}
}

func (i Beqz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("beqz %s, %s", laRegName(i.rj), i.off.text())
}

func (i Beqz) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 21, "beqz offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["beqz"][0] | uint32(i.rj)<<5 | scatterD5k16(off)

	return writeWord(w, word)
}

func (i Beqz) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"beqz",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "off": i.off.val},
	)
}
