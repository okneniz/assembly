package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CtzD - ctz.d rd, rj (2R): rd = the count of trailing zeros of rj.
type CtzD struct {
	base

	rd, rj uint8
}

// CtzD - ctz.d rd, rj.
func (Builder) CtzD(rd, rj Reg) Instr {
	return CtzD{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeCtzD(w uint32, addr uint64) Instr {
	return CtzD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i CtzD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ctz.d %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i CtzD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["ctz.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i CtzD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"ctz.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
