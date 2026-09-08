package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// StgtB - stgt.b rd, rj, rk (DJK): store the low byte of rd with a bounds
// check against rk, trapping outside.
type StgtB struct {
	base

	rd, rj, rk uint8
}

// StgtB - stgt.b rd, rj, rk.
func (Builder) StgtB(rd, rj, rk Reg) Instr {
	return StgtB{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeStgtB(w uint32) Instr {
	return StgtB{
		base: newBase(w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i StgtB) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("stgt.b %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i StgtB) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["stgt.b"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
