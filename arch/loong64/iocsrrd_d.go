package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdD - iocsrrd.d rd, rj (2R): rd = the IOCSR doubleword at rj.
type IocsrrdD struct {
	rd, rj uint8
}

func (i IocsrrdD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrrd.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
