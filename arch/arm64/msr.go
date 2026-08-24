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

func decodeMsr(w uint32, addr uint64) Instr {
	return Msr{
		base:   newBase(addr, w),
		rt:     regNameX(w & 0x1f),
		sysreg: sysRegName(w >> 5 & 0x7fff),
	}
}

func (i Msr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("msr %s, %s", i.sysreg, i.rt)
}

func (i Msr) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, 0xD5100000|regBitsX(i.rt)|invSysRegChecked(i.sysreg)<<5)
}

func (i Msr) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"msr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"System",
		map[string]any{"Rt": i.rt, "sysreg": i.sysreg},
	)
}
