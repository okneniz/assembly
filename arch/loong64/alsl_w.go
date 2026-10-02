package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AlslW - alsl.w rd, rj, rk, shift (DJKUa2): rd = rj + (rk << shift); the
// 1..4 shift is encoded as ui2 = shift - 1 (decode adds 1 back).
type AlslW struct {
	rd, rj, rk uint8
	shift      imm
}

func (i AlslW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"alsl.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		laRegName(i.rk),
		i.shift.text(),
	)
}

func (i AlslW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["alsl.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10 |
		scatterU(i.shift.val-1, 15, 2)

	return writeWord(w, word)
}
