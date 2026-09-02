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

// Encodings of the signed-offset form: the access size is set by rt.
const (
	ldpXEnc uint32 = 0xA9400000 // ldp xt, xt2, [xn, #imm7<<3]
	ldpWEnc uint32 = 0x29400000 // ldp wt, wt2, [xn, #imm7<<2]
)

// Ldp — ldp rt, rt2, [rn, #off]: byte offset, scaling to the
// access size is hidden here. rt/rt2 — x/w registers (register 31 reads
// as zr, the widths must match), rn — x register or SP (register 31 in
// the base reads as sp); off — the signed imm7 range scaled by the
// access size.
func (Builder) Ldp(rt, rt2, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldp", "rt", "x/w register (register 31 in rt reads as zr)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rt2, "Ldp", "rt2", "x/w register (register 31 in rt2 reads as zr)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldp",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireWidth("Ldp", rt, rt2); err != nil {
		return nil, err
	}

	enc, scale := ldpXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = ldpWEnc, 2
	}

	if err := requirePairOff("Ldp", off, scale); err != nil {
		return nil, err
	}

	return Ldp{
		pairBase: newPairBase(rt.name(), rt2.name(), rn.name(), memImm, int64(off), scale, enc),
	}, nil
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
