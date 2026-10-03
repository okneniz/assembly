package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Msr — msr sysreg, rt.
type Msr struct {
	rt, sysreg string
}

// newMsr - the Msr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newMsr(sysreg string, rt Reg) (Msr, error) {
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

	key, err := invSysReg(sysreg)
	if err != nil {
		return Msr{}, fmt.Errorf("arm64.NewMsr: operand sysreg: %w", err)
	}

	return Msr{
		rt:     rt.name(),
		sysreg: sysRegName(key),
	}, nil
}

func (i Msr) Encode(w io.Writer) (int64, error) {
	return writeWord(w, 0xD5100000|regBitsX(i.rt)|invSysRegChecked(i.sysreg)<<5)
}

func (i Msr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("msr %s, %s", i.sysreg, i.rt)
}
