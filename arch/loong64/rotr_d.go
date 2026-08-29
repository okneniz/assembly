package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RotrD - rotr.d rd, rj, rk (3R): rd = rotate_right64(rj, rk & 63).
type RotrD struct {
	base

	rd, rj, rk uint8
}

// NewRotrD - rotr.d rd, rj, rk.
func NewRotrD(rd, rj, rk Reg) Instr {
	return RotrD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeRotrD(w uint32, addr uint64) Instr {
	return RotrD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i RotrD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rotr.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i RotrD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["rotr.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i RotrD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"rotr.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
