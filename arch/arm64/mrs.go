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

// newMrs - the Mrs constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newMrs(b base, rd string, sysreg string) Mrs {
	return Mrs{
		base:   b,
		rd:     rd,
		sysreg: sysreg,
	}
}

func (i Mrs) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mrs %s, %s", i.rd, i.sysreg)
}

func (i Mrs) Encode(w io.Writer) (int64, error) {
	return writeWord(w, 0xD5300000|regBitsX(i.rd)|invSysRegChecked(i.sysreg)<<5)
}

// Mrs — mrs rd, sysreg. rd — only x registers (register 31 reads
// as zr); sysreg — an architectural name from the registry (MIDR_EL1,
// NZCV, ...) or the objdump form S<op0>_<op1>_C<CRn>_C<CRm>_<op2>
// (see invSysReg).
func (Builder) Mrs(rd Reg, sysreg string) (Instr, error) {
	if err := requireClass(
		rd,
		"Mrs",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if _, err := invSysReg(sysreg); err != nil {
		return nil, fmt.Errorf("arm64.NewMrs: operand sysreg: %w", err)
	}

	return newMrs(base{}, rd.name(), sysreg), nil
}

func decodeMrs(w uint32) Instr {
	return newMrs(newBase(w), regNameX(w&0x1f), sysRegName(w>>5&0x7fff))
}
