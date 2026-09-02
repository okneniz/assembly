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

	return Mrs{rd: rd.name(), sysreg: sysreg}, nil
}

func decodeMrs(w uint32, addr uint64) Instr {
	return Mrs{
		base:   newBase(addr, w),
		rd:     regNameX(w & 0x1f),
		sysreg: sysRegName(w >> 5 & 0x7fff),
	}
}

func (i Mrs) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mrs %s, %s", i.rd, i.sysreg)
}

func (i Mrs) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, 0xD5300000|regBitsX(i.rd)|invSysRegChecked(i.sysreg)<<5)
}

func (i Mrs) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"mrs",
		i.ObjDump(disasm.DefaultViewCtx()),
		"System",
		map[string]any{"Rd": i.rd, "sysreg": i.sysreg},
	)
}
