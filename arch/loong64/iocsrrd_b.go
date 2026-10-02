package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrrdB - iocsrrd.b rd, rj (2R): rd = the IOCSR byte at rj.
type IocsrrdB struct {
	rd, rj uint8
}

func (i IocsrrdB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrrd.b %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrrdB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrrd.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
