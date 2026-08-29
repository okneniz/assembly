package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWBW - crc.w.b.w rd, rj, rk (DJK): CRC32 of rj (byte) into rk's
// accumulator.
type CrcWBW struct {
	base

	rd, rj, rk uint8
}

// NewCrcWBW - crc.w.b.w rd, rj, rk.
func NewCrcWBW(rd, rj, rk Reg) Instr {
	return CrcWBW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrcWBW(w uint32, addr uint64) Instr {
	return CrcWBW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrcWBW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.b.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWBW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crc.w.b.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrcWBW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crc.w.b.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
