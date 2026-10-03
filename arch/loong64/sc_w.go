package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// ScW - sc.w rd, rj, offs (DJSk14): MEM[rj + offs] = low32(rd)
// conditionally, rd = success (store-conditional 32 bits; offs - a byte
// offset like ll.w).
type ScW struct {
	rd, rj uint8
	off    imm
}

func (i ScW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["sc.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 | scatterS(i.off.val>>2, 10, 14)

	return writeWord(w, word)
}

func (i ScW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sc.w %s, %s, %s", laRegName(i.rd), laRegName(i.rj), i.off.text())
}
