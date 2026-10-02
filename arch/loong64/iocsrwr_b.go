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

func decodeIocsrwrB(w uint32) Instr {
	return IocsrwrB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i IocsrwrB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrwr.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
