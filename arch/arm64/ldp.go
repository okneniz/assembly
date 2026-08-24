package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldp — ldp rt, rt2, [rn{, #imm7<<scale}{!}].
type Ldp struct {
	base
	pairBase
}

func decodeLdpOf(enc uint32, scale uint32, x64 bool, rtKind string) func(uint32, uint64) Instr {
	kind := rtKind
	if kind == "" {
		if x64 {
			kind = "x"
		} else {
			kind = "w"
		}
	}

	return func(w uint32, addr uint64) Instr {
		rt, rt2, rn, k, off, load := pairDecode(w, scale, kind)
		if !load {
			return Stp{
				base:     newBase(addr, w),
				pairBase: newPairBase(rt, rt2, rn, k, off, scale, enc&^1<<22),
			}
		}

		return Ldp{
			base:     newBase(addr, w),
			pairBase: newPairBase(rt, rt2, rn, k, off, scale, enc|1<<22),
		}
	}
}

func (i Ldp) ObjDump(_ disasm.ViewCtx) string {
	return "ldp " + i.pairText()
}

func (i Ldp) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.pairWrite(w, "ldp")
}

func (i Ldp) MarshalJSON() ([]byte, error) {
	return i.marshal("ldp", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rt": i.rt, "Rt2": i.rt2, "Rn": i.rn})
}
