package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// RdtimehW - rdtimeh.w rd, rj (2R): rd = high 32 bits of the stable counter + rj.
type RdtimehW struct {
	rd, rj uint8
}

func (i RdtimehW) ObjDump(_ disasm.ViewCtx) string {
	if i.rj == 0 {
		return "rdcntvh.w " + laRegName(i.rd)
	}

	return fmt.Sprintf("rdtimeh.w %s, %s", laRegName(i.rd), laRegName(i.rj))
}

func (i RdtimehW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["rdtimeh.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5

	return writeWord(w, word)
}
