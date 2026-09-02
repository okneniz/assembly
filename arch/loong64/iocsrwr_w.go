package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrwrW - iocsrwr.w rd, rj (2R): the IOCSR word at rj = rd.
type IocsrwrW struct {
	base

	rd, rj uint8
}

// IocsrwrW - iocsrwr.w rd, rj.
func (Builder) IocsrwrW(rd, rj Reg) Instr {
	return IocsrwrW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrwrW(w uint32, addr uint64) Instr {
	return IocsrwrW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrwrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrwr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrwrW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrwr.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
