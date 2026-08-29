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

// NewAsrtleD - asrtle.d rj, rk.
func NewAsrtleD(rj, rk Reg) Instr {
	return AsrtleD{
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeAsrtleD(w uint32, addr uint64) Instr {
	return AsrtleD{
		base: newBase(addr, w),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i AsrtleD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("asrtle.d %s, %s", laRegName(i.rj), laRegName(i.rk))
}

func (i AsrtleD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["asrtle.d"][0] |
		uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i AsrtleD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"asrtle.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
