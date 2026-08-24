package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stxrb — stxrb rs, rt, [rn].
type Stxrb struct {
	base
	excl

	enc uint32
}

func decodeStxrbOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stxrb{
			base: newBase(addr, w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}
	}
}

func (i Stxrb) ObjDump(_ disasm.ViewCtx) string {
	return "stxrb " + i.exText()
}

func (i Stxrb) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.exWrite(w, i.enc, "stxrb")
}

func (i Stxrb) MarshalJSON() ([]byte, error) {
	return i.marshal("stxrb", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rs": i.rs, "Rt": i.rt, "Rn": i.rn})
}
