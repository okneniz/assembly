package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// CrcWWW - crc.w.w.w rd, rj, rk (DJK): CRC32 of rj (word) into rk's
// accumulator.
type CrcWWW struct {
	rd, rj, rk uint8
}

func (i CrcWWW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("crc.w.w.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), laRegName(i.rk))
}

func (i CrcWWW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["crc.w.w.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10

	return writeWord(w, word)
}
