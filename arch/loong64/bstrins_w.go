package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrinsW - bstrins.w rd, rj, msb, lsb (DJUk5Um5): insert the rj[msb:lsb] field into rd[msb:lsb].
type BstrinsW struct {
	rd, rj uint8
	msb    imm
	lsb    imm
}

func (i BstrinsW) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bstrins.w"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 5) | scatterU(i.lsb.val, 10, 5)

	return writeWord(w, word)
}

func (i BstrinsW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrins.w %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}
