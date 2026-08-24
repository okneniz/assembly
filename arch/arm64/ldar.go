package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldar — ldar rt, [rn].
type Ldar struct {
	base
	atomic

	enc uint32
}

func decodeLdarOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Ldar{
			base:   newBase(addr, w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}
	}
}

func (i Ldar) ObjDump(_ disasm.ViewCtx) string {
	return "ldar " + i.atText()
}

func (i Ldar) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.atWrite(w, i.enc, "ldar")
}

func (i Ldar) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldar",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
