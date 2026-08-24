package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldarb — ldarb rt, [rn].
type Ldarb struct {
	base
	atomic

	enc uint32
}

func decodeLdarbOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Ldarb{
			base:   newBase(addr, w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}
	}
}

func (i Ldarb) ObjDump(_ disasm.ViewCtx) string {
	return "ldarb " + i.atText()
}

func (i Ldarb) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.atWrite(w, i.enc, "ldarb")
}

func (i Ldarb) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldarb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
