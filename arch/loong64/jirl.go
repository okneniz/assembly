package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Jirl - jirl rd, rj, offs (2RI16): rd = pc + 4, jump to rj + offs (offs
// is word-scaled, relative to the rj VALUE - unknown at decode time, so
// the raw byte offset is stored, like riscv jalr).
type Jirl struct {
	base

	rd, rj uint8
	off    imm
}

// Jirl - jirl rd, rj, offs (the byte offset from rj).
func (Builder) Jirl(rd, rj Reg, off int64) Instr {
	return Jirl{
		rd:  rd.Num(),
		rj:  rj.Num(),
		off: immNum(off),
	}
}

func decodeJirl(w uint32, addr uint64) Instr {
	return Jirl{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		off:  immNum(sField(w, 10, 16) << 2),
	}
}

func (i Jirl) ObjDump(_ disasm.ViewCtx) string {
	if i.rd == 0 && i.off.val == 0 {
		if i.rj == 1 {
			return "ret"
		}

		return "jr " + laRegName(i.rj)
	}

	return fmt.Sprintf("jirl %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i Jirl) Encode(w io.Writer, _ uint64) (int64, error) {
	off, err := encPs2(i.off.val, 16, "jirl offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["jirl"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}

func (i Jirl) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"jirl",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "off": i.off.val},
	)
}
