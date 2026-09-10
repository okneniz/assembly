package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Mrs — mrs rd, sysreg.
type Mrs struct {
	base

	rd, sysreg string
}

// newMrs - the Mrs constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newMrs(b base, rd Reg, sysreg string) (Mrs, error) {
	err := requireClass(
		rd,
		"Mrs",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Mrs{}, err
	}

	if _, err := invSysReg(sysreg); err != nil {
		return Mrs{}, fmt.Errorf("arm64.NewMrs: operand sysreg: %w", err)
	}

	return Mrs{
		base:   b,
		rd:     rd.name(),
		sysreg: sysreg,
	}, nil
}

func (i Mrs) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mrs %s, %s", i.rd, i.sysreg)
}

func (i Mrs) Encode(w io.Writer) (int64, error) {
	return writeWord(w, 0xD5300000|regBitsX(i.rd)|invSysRegChecked(i.sysreg)<<5)
}

func (Builder) Mrs(rd Reg, sysreg string) (Instr, error) {
	return newMrs(base{}, rd, sysreg)
}

func decodeMrs(w uint32) (Instr, error) {
	in, err := newMrs(newBase(w), gprOf(w&0x1f, true), sysRegName(w>>5&0x7fff))
	if err != nil {
		return nil, err
	}

	return in, nil
}
