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

// Mv - mv rd, rs2 (the c.mv halfword).
func (Builder) Mv(rd, rs2 Reg) Instr {
	h := uint32(0x8002) | uint32(rd.Num())<<7 | uint32(rs2.Num())<<2

	return Mv{
		base: newHalfBase(h),
		rd:   rd.name(),
		rs2:  rs2.name(),
	}
}

// cMv - compressed forms (c.mv): base - halfword, length 2.
func cMv(h uint32, rd, rs2 string) Mv {
	return Mv{
		base: newHalfBase(h),
		rd:   rd,
		rs2:  rs2,
	}
}

func (i Mv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mv %s, %s", i.rd, i.rs2)
}

func (i Mv) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeHalf(w, uint16(i.raw))
}
