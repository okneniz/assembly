package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdH - iocsrrd.h rd, rj (2R): rd = the IOCSR halfword at rj.
type IocsrrdH struct {
	base

	rd, rj uint8
}

// NewIocsrrdH - iocsrrd.h rd, rj.
func NewIocsrrdH(rd, rj Reg) Instr {
	return IocsrrdH{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrrdH(w uint32, addr uint64) Instr {
	return IocsrrdH{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrrdH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdH) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrrd.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrrdH) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrrd.h",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
