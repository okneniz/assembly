package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWBW - crc.w.b.w rd, rj, rk (DJK): CRC32 of rj (byte) into rk's
// accumulator.
type CrcWBW struct {
	rd, rj, rk uint8
}

func (i CrcWBW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.b.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWBW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crc.w.b.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
