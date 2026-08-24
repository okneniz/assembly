package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsb — ldrsb rt, [rn, #imm12] (Rt — X; the ldrsbr W form is not in the table).
type Ldrsb struct {
	base

	rt, rn string
	off    int64
}

const ldrsbEnc uint32 = 0x39800000

func decodeLdrsb(w uint32, addr uint64) Instr {
	return Ldrsb{
		base: newBase(addr, w),
		rt:   regNameX(w & 0x1f),
		rn:   regNameXSP(w >> 5 & 0x1f),
		off:  int64(w >> 10 & 0xfff),
	}
}

func (i Ldrsb) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldrsb %s, [%s]", i.rt, i.rn)
	}

	return fmt.Sprintf("ldrsb %s, [%s, #0x%x]", i.rt, i.rn, i.off)
}

func (i Ldrsb) Encode(w io.Writer, pc uint64) (int64, error) {
	return lsSignedWrite(w, ldrsbEnc, i.rt, i.rn, i.off, "ldrsb")
}

func (i Ldrsb) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldrsb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
