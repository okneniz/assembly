package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsh — ldrsh rt, [rn, #imm12<<1].
type Ldrsh struct {
	base

	rt, rn string
	off    int64
}

const ldrshEnc uint32 = 0x79800000

func decodeLdrsh(w uint32, addr uint64) Instr {
	return Ldrsh{
		base: newBase(addr, w),
		rt:   regNameX(w & 0x1f),
		rn:   regNameXSP(w >> 5 & 0x1f),
		off:  int64(w>>10&0xfff) << 1,
	}
}

func (i Ldrsh) ObjDump(_ disasm.ViewCtx) string {
	if i.off == 0 {
		return fmt.Sprintf("ldrsh %s, [%s]", i.rt, i.rn)
	}

	return fmt.Sprintf("ldrsh %s, [%s, #0x%x]", i.rt, i.rn, i.off)
}

func (i Ldrsh) Encode(w io.Writer, pc uint64) (int64, error) {
	return lsSignedWrite(w, ldrshEnc, i.rt, i.rn, i.off, "ldrsh")
}

func (i Ldrsh) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldrsh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
