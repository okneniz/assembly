package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// BstrinsD - bstrins.d rd, rj, msb, lsb (DJUk6Um6): insert the rj[msb:lsb] field into rd[msb:lsb].
type BstrinsD struct {
	rd, rj uint8
	msb    imm
	lsb    imm
}

func (i BstrinsD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf(
		"bstrins.d %s, %s, %s, %s",
		laRegName(i.rd),
		laRegName(i.rj),
		i.msb.text(),
		i.lsb.text(),
	)
}

func (i BstrinsD) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["bstrins.d"][0] |
		uint32(i.rd) | uint32(i.rj)<<5 |
		scatterU(i.msb.val, 16, 6) | scatterU(i.lsb.val, 10, 6)

	return writeWord(w, word)
}
