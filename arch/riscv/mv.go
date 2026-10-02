package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mv - c.mv: displayed in the pseudo-form mv (a 32-bit add with rs1=zero
// would print "add rd, zero, rs", while c.mv prints exactly "mv"). Decode side:
// mv expansion during assembly is in asm/riscv/pseudo.
type Mv struct {
	half uint32

	rd, rs2 string
}

// cMv - compressed forms (c.mv): the halfword is the encoding source.
func cMv(h uint32, rd, rs2 string) Mv {
	return Mv{
		half: h,
		rd:   rd,
		rs2:  rs2,
	}
}

func (i Mv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mv %s, %s", i.rd, i.rs2)
}

func (i Mv) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeHalf(w, uint16(i.half))
}
