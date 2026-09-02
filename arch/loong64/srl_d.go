package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SrlD - srl.d rd, rj, rk (3R): rd = rj >>l (rk & 63).
type SrlD struct {
	base

	rd, rj, rk uint8
}

// SrlD - srl.d rd, rj, rk.
func (Builder) SrlD(rd, rj, rk Reg) Instr {
	return SrlD{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeSrlD(w uint32, addr uint64) Instr {
	return SrlD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i SrlD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srl.d %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i SrlD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["srl.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i SrlD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srl.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
