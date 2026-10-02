package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AsrtleD - asrtle.d rj, rk (JK): raise a bounds trap unless rj <= rk.
type AsrtleD struct {
	base

	rj, rk uint8
}

func (i AsrtleD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asrtle.d %s, %s", laRegName(i.rj), laRegName(i.rk))
}

func (i AsrtleD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["asrtle.d"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
