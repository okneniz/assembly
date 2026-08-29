package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdW - iocsrrd.w rd, rj (2R): rd = the IOCSR word at rj.
type IocsrrdW struct {
	base

	rd, rj uint8
}

// NewIocsrrdW - iocsrrd.w rd, rj.
func NewIocsrrdW(rd, rj Reg) Instr {
	return IocsrrdW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrrdW(w uint32, addr uint64) Instr {
	return IocsrrdW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrrdW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrrd.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrrdW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrrd.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
