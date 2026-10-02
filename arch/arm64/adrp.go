package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Adrp — adrp rd, #imm21 (the absolute-page annotation is computed at
// print time from the view-context address).
type Adrp struct {
	rd  string
	off int64
}

// newAdrp - the Adrp constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAdrp(rd Reg, off int64) (Adrp, error) {
	err := requireClass(
		rd,
		"Adrp",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Adrp{}, err
	}

	if off < -(1<<20) || off >= 1<<20 {
		return Adrp{}, fmt.Errorf(
			"arm64.NewAdrp: operand off: %d is out of the imm21 range (-0x100000..0xfffff)",
			off,
		)
	}

	return Adrp{
		rd:  rd.name(),
		off: off,
	}, nil
}

func (i Adrp) ObjDump(ctx disasm.ViewCtx) string {
	page := int64(ctx.Addr())&^int64(0xFFF) + i.off<<12
	return fmt.Sprintf("adrp %s, %d ; 0x%x", i.rd, i.off, page)
}

func (i Adrp) Encode(w io.Writer) (int64, error) {
	return writeWord(w, 0x90000000|regBitsX(i.rd)|uint32(i.off&3)<<29|uint32(i.off>>2&0x7ffff)<<5)
}
