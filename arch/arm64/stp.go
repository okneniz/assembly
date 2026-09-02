package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stp — stp rt, rt2, [rn{, #imm7<<scale}{!}] (the same encoding, L=0).
type Stp struct {
	base
	pairBase
}

// Encodings of the signed-offset form: the access size is set by rt.
const (
	stpXEnc uint32 = 0xA9000000 // stp xt, xt2, [xn, #imm7<<3]
	stpWEnc uint32 = 0x29000000 // stp wt, wt2, [xn, #imm7<<2]
)

// Stp — stp rt, rt2, [rn, #off]: byte offset, scaling to the
// access size is hidden here; the operand constraints are those of Ldp.
func (Builder) Stp(rt, rt2, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Stp", "rt", "x/w register (register 31 in rt reads as zr)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rt2, "Stp", "rt2", "x/w register (register 31 in rt2 reads as zr)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Stp",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireWidth("Stp", rt, rt2); err != nil {
		return nil, err
	}

	enc, scale := stpXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = stpWEnc, 2
	}

	if err := requirePairOff("Stp", off, scale); err != nil {
		return nil, err
	}

	return Stp{
		pairBase: newPairBase(rt.name(), rt2.name(), rn.name(), memImm, int64(off), scale, enc),
	}, nil
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
