package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrpickD - bstrpick.d rd, rj, msb, lsb (DJUk6Um6): extract the rj[msb:lsb] field into rd, zero-extended.
type BstrpickD struct {
	base

	rd, rj uint8
	msb    imm
	lsb    imm
}

func (i BstrpickD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrpick.d %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrpickD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bstrpick.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 6) | scatterU(i.lsb.val, 10, 6)

	return writeWord(w, word)
}
