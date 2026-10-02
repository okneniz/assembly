package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdHu - ld.hu rd, rj, si12 (2RI12): rd = the halfword at rj + si12, zero-extended.
type LdHu struct {
	rd, rj uint8
	off    imm
}

func (i LdHu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.hu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdHu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ld.hu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}
