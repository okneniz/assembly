package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrwrD - iocsrwr.d rd, rj (2R): the IOCSR doubleword at rj = rd.
type IocsrwrD struct {
	base

	rd, rj uint8
}

// NewIocsrwrD - iocsrwr.d rd, rj.
func NewIocsrwrD(rd, rj Reg) Instr {
	return IocsrwrD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrwrD(w uint32, addr uint64) Instr {
	return IocsrwrD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrwrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrwr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrwrD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrwr.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
