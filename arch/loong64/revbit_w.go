package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RevbitW - bitrev.w rd, rj (2R): rd = low32(rj) with all 32 bits reversed.
type RevbitW struct {
	base

	rd, rj uint8
}

// RevbitW - bitrev.w rd, rj.
func (Builder) RevbitW(rd, rj Reg) Instr {
	return RevbitW{
		rd: rd.Num(),
		rj: rj.Num(),
	}
}

func decodeRevbitW(w uint32, addr uint64) Instr {
	return RevbitW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
	}
}

func (i RevbitW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bitrev.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RevbitW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["bitrev.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}

func (i RevbitW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bitrev.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj)},
	)
}
