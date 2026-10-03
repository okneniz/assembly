package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BytepickW - bytepick.w rd, rj, rk, sel (DJKUa2): pick 4 bytes out of the
// {rj, rk} concatenation, byte-indexed by sel.
type BytepickW struct {
	rd, rj, rk uint8
	sel        imm
}

func (i BytepickW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bytepick.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | uint32(i.rk)<<10 |
		scatterU(i.sel.val, 15, 2)

	return writeWord(w, word)
}

func (i BytepickW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bytepick.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		laRegName(i.rk),
		i.sel.text(),
	)
}
