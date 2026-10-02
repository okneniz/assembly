package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AlslWu - alsl.wu rd, rj, rk, shift (DJKUa2): rd = zero-extend32(rj +
// (rk << shift)); the 1..4 shift is encoded as ui2 = shift - 1.
type AlslWu struct {
	rd, rj, rk uint8
	shift      imm
}

func (i AlslWu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"alsl.wu %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		laRegName(i.rk),
		i.shift.text(),
	)
}

func (i AlslWu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["alsl.wu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10 |
		scatterU(i.shift.val-1, 15, 2)

	return writeWord(w, word)
}
