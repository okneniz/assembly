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

func (i CrcWDW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.d.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWDW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crc.w.d.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
