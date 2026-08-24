package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stp - stp rt, rt2, [rn{, #imm7<<scale}{!}] (the same encoding, L=0).
type Stp struct {
	base
	pairBase
}

func (i Stp) ObjDump(_ disasm.ViewCtx) string {
	return "stp " + i.pairText()
}

func (i Stp) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.pairWrite(w, "stp")
}

func (i Stp) MarshalJSON() ([]byte, error) {
	return i.marshal("stp", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rt": i.rt, "Rt2": i.rt2, "Rn": i.rn})
}
