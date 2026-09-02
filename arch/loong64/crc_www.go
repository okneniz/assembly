package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWWW - crc.w.w.w rd, rj, rk (DJK): CRC32 of rj (word) into rk's
// accumulator.
type CrcWWW struct {
	base

	rd, rj, rk uint8
}

// CrcWWW - crc.w.w.w rd, rj, rk.
func (Builder) CrcWWW(rd, rj, rk Reg) Instr {
	return CrcWWW{
		rd: rd.Num(),
		rj: rj.Num(),
		rk: rk.Num(),
	}
}

func decodeCrcWWW(w uint32, addr uint64) Instr {
	return CrcWWW{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		rj:   uint8(w >> 5 & 0x1f),
		rk:   uint8(w >> 10 & 0x1f),
	}
}

func (i CrcWWW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.w.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWWW) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["crc.w.w.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrcWWW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"crc.w.w.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "rj": laRegName(i.rj), "rk": laRegName(i.rk)},
	)
}
