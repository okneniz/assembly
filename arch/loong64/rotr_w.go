package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RotrW - rotr.w rd, rj, rk (3R): rd = sign32(rotate_right32(rj, rk & 31)).
type RotrW struct {
	base

	rd, rj, rk uint8
}

// RotrW - rotr.w rd, rj, rk.
func (Builder) RotrW(rd, rj, rk Reg) Instr {
	return RotrW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeRotrW(w uint32, addr uint64) Instr {
	return RotrW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i RotrW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rotr.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i RotrW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["rotr.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i RotrW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"rotr.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
