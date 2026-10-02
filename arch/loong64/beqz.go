package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Beqz - beqz rj, offs (1RI21, the immediate split d5k16): if rj == 0,
// jump to pc + offs (word-scaled). The decoded form stores the absolute
// target.
type Beqz struct {
	base

	rj  uint8
	off imm
}

func (i Beqz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("beqz %s, %s", laRegName(i.rj), i.off.text())
}

func (i Beqz) Encode(w io.Writer) (int64, error) {
	off, err := encPs2(i.off.val, 21, "beqz offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["beqz"][0] | uint32(i.rj)<<5 | scatterD5k16(off)

	return writeWord(w, word)
}
