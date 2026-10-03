package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bltu - bltu rj, rd, offs (2RI16, the manual order swaps the registers):
// if the unsigned rj < rd, jump to pc + offs (offs is word-scaled). The
// decoded form stores the absolute target.
type Bltu struct {
	rd, rj uint8
	off    imm
}

func (i Bltu) Encode(w io.Writer) (int64, error) {
	off, err := encPs2(i.off.val, 16, "bltu offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bltu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(off, 10, 16)

	return writeWord(w, word)
}

func (i Bltu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bltu %s, %s, %s", laRegName(i.rj), laRegName(i.rd), i.off.text())
}
