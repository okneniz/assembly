package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrwrB - iocsrwr.b rd, rj (2R): the IOCSR byte at rj = rd.
type IocsrwrB struct {
	base

	rd, rj uint8
}

// NewIocsrwrB - iocsrwr.b rd, rj.
func NewIocsrwrB(rd, rj Reg) Instr {
	return IocsrwrB{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrwrB(w uint32, addr uint64) Instr {
	return IocsrwrB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrwrB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrwr.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrwrB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrwr.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
