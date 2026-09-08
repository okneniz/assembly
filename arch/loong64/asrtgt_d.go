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

// AsrtgtD - asrtgt.d rj, rk.
func (Builder) AsrtgtD(rj, rk Reg) Instr {
	return AsrtgtD{
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAsrtgtD(w uint32) Instr {
	return AsrtgtD{
		base: newBase(w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i AsrtgtD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asrtgt.d %s, %s", laRegName(i.rj), laRegName(i.rk))
}

func (i AsrtgtD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["asrtgt.d"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
