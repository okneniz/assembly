package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrpickW - bstrpick.w rd, rj, msb, lsb (DJUk5Um5): extract the rj[msb:lsb] field into rd, zero-extended.
type BstrpickW struct {
	rd, rj uint8
	msb    imm
	lsb    imm
}

func (i BstrpickW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrpick.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrpickW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bstrpick.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 5) | scatterU(i.lsb.val, 10, 5)

	return writeWord(w, word)
}
