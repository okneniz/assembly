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
	base

	rd, rs2 string
}

// cMv - compressed forms (c.mv): base - halfword, length 2.
func cMv(h uint32, addr uint64, rd, rs2 string) Mv {
	return Mv{
		base: newHalfBase(h, addr),
		rd:   rd,
		rs2:  rs2,
	}
}

func (i Mv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mv %s, %s", i.rd, i.rs2)
}

func (i Mv) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeHalf(w, uint16(i.raw))
}

func (i Mv) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"mv",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Pseudo",
		map[string]any{"rd": i.rd, "rs2": i.rs2},
	)
}
