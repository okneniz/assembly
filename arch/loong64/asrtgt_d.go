package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AsrtgtD - asrtgt.d rj, rk (JK): raise a bounds trap unless rj > rk.
type AsrtgtD struct {
	base

	rj, rk uint8
}

func (i AsrtgtD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asrtgt.d %s, %s", laRegName(i.rj), laRegName(i.rk))
}

func (i AsrtgtD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["asrtgt.d"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
