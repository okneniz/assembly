package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWDW - crc.w.d.w rd, rj, rk (DJK): CRC32 of rj (double) into rk's
// accumulator.
type CrcWDW struct {
	base

	rd, rj, rk uint8
}

// NewCrcWDW - crc.w.d.w rd, rj, rk.
func NewCrcWDW(rd, rj, rk Reg) Instr {
	return CrcWDW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrcWDW(w uint32, addr uint64) Instr {
	return CrcWDW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrcWDW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.d.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWDW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crc.w.d.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrcWDW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crc.w.d.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
