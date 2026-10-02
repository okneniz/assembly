package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrwrH - iocsrwr.h rd, rj (2R): the IOCSR halfword at rj = rd.
type IocsrwrH struct {
	rd, rj uint8
}

func (i IocsrwrH) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.h %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i IocsrwrH) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrwr.h"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
