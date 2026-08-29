package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdD - iocsrrd.d rd, rj (2R): rd = the IOCSR doubleword at rj.
type IocsrrdD struct {
	base

	rd, rj uint8
}

// NewIocsrrdD - iocsrrd.d rd, rj.
func NewIocsrrdD(rd, rj Reg) Instr {
	return IocsrrdD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrrdD(w uint32, addr uint64) Instr {
	return IocsrrdD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrrdD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrrd.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrrdD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrrd.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
