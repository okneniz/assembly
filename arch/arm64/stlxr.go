package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlxr — stlxr rs, rt, [rn].
type Stlxr struct {
	base
	excl

	enc uint32
}

func decodeStlxrOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stlxr{
			base: newBase(addr, w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}
	}
}

func (i Stlxr) ObjDump(_ disasm.ViewCtx) string {
	return "stlxr " + i.exText()
}

func (i Stlxr) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.exWrite(w, i.enc, "stlxr")
}

func (i Stlxr) MarshalJSON() ([]byte, error) {
	return i.marshal("stlxr", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rs": i.rs, "Rt": i.rt, "Rn": i.rn})
}
