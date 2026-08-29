package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdB - iocsrrd.b rd, rj (2R): rd = the IOCSR byte at rj.
type IocsrrdB struct {
	base

	rd, rj uint8
}

// NewIocsrrdB - iocsrrd.b rd, rj.
func NewIocsrrdB(rd, rj Reg) Instr {
	return IocsrrdB{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeIocsrrdB(w uint32, addr uint64) Instr {
	return IocsrrdB{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrrdB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdB) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["iocsrrd.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrrdB) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"iocsrrd.b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Privileged",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
