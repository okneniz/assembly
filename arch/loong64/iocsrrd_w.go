package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdW - iocsrrd.w rd, rj (2R): rd = the IOCSR word at rj.
type IocsrrdW struct {
	rd, rj uint8
}

func (i IocsrrdW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrrd.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
