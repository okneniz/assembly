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

// NewAsrtgtD - asrtgt.d rj, rk.
func NewAsrtgtD(rj, rk Reg) Instr {
	return AsrtgtD{
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAsrtgtD(w uint32, addr uint64) Instr {
	return AsrtgtD{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i AsrtgtD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asrtgt.d %s, %s", laRegName(i.rj), laRegName(i.rk))
}

func (i AsrtgtD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["asrtgt.d"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i AsrtgtD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"asrtgt.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
