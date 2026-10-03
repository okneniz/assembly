package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdH - iocsrrd.h rd, rj (2R): rd = the IOCSR halfword at rj.
type IocsrrdH struct {
	rd, rj uint8
}

func (i IocsrrdH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrrd.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrrdH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.h %s, %s", laRegName(i.rd), laRegName(i.rj))
}
