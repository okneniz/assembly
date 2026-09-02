package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrwrH - iocsrwr.h rd, rj (2R): the IOCSR halfword at rj = rd.
type IocsrwrH struct {
	base

	rd, rj uint8
}

// IocsrwrH - iocsrwr.h rd, rj.
func (Builder) IocsrwrH(rd, rj Reg) Instr {
	return IocsrwrH{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrwrH(w uint32, addr uint64) Instr {
	return IocsrwrH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrwrH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrwr.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrwrH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrwr.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
