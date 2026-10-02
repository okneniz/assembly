package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bgeu - bgeu rj, rd, offs (2RI16, the manual order swaps the registers):
// if the unsigned rj >= rd, jump to pc + offs (offs is word-scaled). The
// decoded form stores the absolute target.
type Bgeu struct {
	base

	rd, rj uint8
	off    imm
}

func (i Bgeu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bgeu %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}

func (i Bgeu) Encode(w io.Writer) (int64, error) {
	off, err := encPs2(i.off.val, 16, "bgeu offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bgeu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}
