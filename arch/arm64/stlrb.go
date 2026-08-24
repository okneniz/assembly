package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlrb — stlrb rt, [rn].
type Stlrb struct {
	base
	atomic

	enc uint32
}

func decodeStlrbOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stlrb{
			base:   newBase(addr, w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}
	}
}

func (i Stlrb) ObjDump(_ disasm.ViewCtx) string {
	return "stlrb " + i.atText()
}

func (i Stlrb) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.atWrite(w, i.enc, "stlrb")
}

func (i Stlrb) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"stlrb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
