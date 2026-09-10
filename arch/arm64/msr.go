package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Msr — msr sysreg, rt.
type Msr struct {
	base

	rt, sysreg string
}

// newMsr - the Msr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newMsr(b base, sysreg string, rt Reg) (Msr, error) {
	err := requireClass(
		rt,
		"Msr",
		"rt",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Msr{}, err
	}

	if _, err := invSysReg(sysreg); err != nil {
		return Msr{}, fmt.Errorf("arm64.NewMsr: operand sysreg: %w", err)
	}

	return Msr{
		base:   b,
		rt:     rt.name(),
		sysreg: sysreg,
	}, nil
}

func (i Msr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("msr %s, %s", i.sysreg, i.rt)
}

func (i Msr) Encode(w io.Writer) (int64, error) {
	return writeWord(w, 0xD5100000|regBitsX(i.rt)|invSysRegChecked(i.sysreg)<<5)
}

func (Builder) Msr(sysreg string, rt Reg) (Instr, error) {
	return newMsr(base{}, sysreg, rt)
}

func decodeMsr(w uint32) (Instr, error) {
	in, err := newMsr(
		newBase(w),
		sysRegName(w>>5&0x7fff),
		gprOf(w&0x1f, true),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
