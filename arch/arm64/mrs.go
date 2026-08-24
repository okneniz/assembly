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
