package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWHW - crc.w.h.w rd, rj, rk (DJK): CRC32 of rj (half) into rk's
// accumulator.
type CrcWHW struct {
	rd, rj, rk uint8
}

func (i CrcWHW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crc.w.h.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}

func (i CrcWHW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.h.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}
