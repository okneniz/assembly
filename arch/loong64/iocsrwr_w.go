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

func (i IocsrwrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrwr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
