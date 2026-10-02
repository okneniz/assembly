package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LdBu - ld.bu rd, rj, si12 (2RI12): rd = the byte at rj + si12, zero-extended.
type LdBu struct {
	rd, rj uint8
	off    imm
}

func (i LdBu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ld.bu %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}

func (i LdBu) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["ld.bu"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val, 10, 12)

	return writeWord(w, word)
}
