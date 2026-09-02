package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsw — ldrsw ... (see lsBase for the addressing kinds).
type Ldrsw struct {
	base
	lsBase
}

const ldrswEnc uint32 = 0xB9800000 // ldrsw xt, [xn, #imm12<<2]

// Ldrsw — ldrsw rt, [rn, #off]: sign-extending word load, rt — x
// register only (register 31 reads as xzr), rn — x register or SP
// (register 31 in the base reads as sp); the offset is an imm12 scaled
// by 4 (0..0xffc, alignment 4).
func (Builder) Ldrsw(rt, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldrsw", "rt", "x register (register 31 in rt reads as xzr)",
		classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldrsw",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireOff("Ldrsw", off, 2); err != nil {
		return nil, err
	}

	return Ldrsw{
		lsBase: newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, ldrswEnc, "", "", 0),
	}, nil
}

func (i Ldrsw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldrsw %s, %s", i.rt, i.lsText())
}

func (i Ldrsw) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.lsWrite(w, pc, "ldrsw")
}

func (i Ldrsw) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldrsw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}

func decodeLdrswOf(enc uint32, kind memKind) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		rt := regNameX(w & 0x1f)
		rn := regNameXSP(w >> 5 & 0x1f)
		var off int64
		var tgt uint64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			tgt = addr + uint64(signExtendN(w>>5&0x7ffff, 19))*4
		case memRegOff:
			rm = regNameX(w >> 16 & 0x1f)
			option = lsOptName(w >> 13 & 7)
			sBit := w>>12&1 == 1
			scale := w >> 30 & 3
			switch {
			case option == "lsl" && sBit && scale > 0:
				shiftAmt = scale
			case option == "lsl":
				option = "" // [rn, rm] without extension
			case sBit && scale > 0:
				shiftAmt = scale
			}
		case memUnscaled, memPost, memPre:
			off = signExtendN(w>>12&0x1ff, 9)
		}

		return Ldrsw{
			base:   newBase(addr, w),
			lsBase: newLsBase(rt, rn, kind, off, tgt, enc, rm, option, shiftAmt),
		}
	}
}
