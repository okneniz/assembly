package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// IocsrwrD - iocsrwr.d rd, rj (2R): the IOCSR doubleword at rj = rd.
type IocsrwrD struct {
	rd, rj uint8
}

func (i IocsrwrD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["iocsrwr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i IocsrwrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("iocsrwr.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}
