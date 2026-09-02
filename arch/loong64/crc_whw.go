package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWHW - crc.w.h.w rd, rj, rk (DJK): CRC32 of rj (half) into rk's
// accumulator.
type CrcWHW struct {
	base

	rd, rj, rk uint8
}

// CrcWHW - crc.w.h.w rd, rj, rk.
func (Builder) CrcWHW(rd, rj, rk Reg) Instr {
	return CrcWHW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrcWHW(w uint32, addr uint64) Instr {
	return CrcWHW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrcWHW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.h.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWHW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crc.w.h.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrcWHW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crc.w.h.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
