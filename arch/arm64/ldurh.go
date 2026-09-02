package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldurh — ldurh ... (see lsBase for the addressing kinds).
type Ldurh struct {
	base
	lsBase
}

const ldurhEnc uint32 = 0x78400000 // ldurh wt, [xn, #±imm9]

// Ldurh — ldurh rt, [rn, #off]: the unscaled form, halfword
// access, rt — w register only (register 31 reads as wzr), rn — x
// register or SP (register 31 in the base reads as sp); the offset is a
// signed imm9 (-0x100..0xff, any alignment).
func (Builder) Ldurh(rt, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldurh", "rt", "w register (register 31 in rt reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldurh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireUnscaledOff("Ldurh", off); err != nil {
		return nil, err
	}

	return Ldurh{
		lsBase: newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, ldurhEnc, "", "", 0),
	}, nil
}

func (i Ldurh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldurh %s, %s", i.rt, i.lsText())
}

func (i Ldurh) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.lsWrite(w, pc, "ldurh")
}

func (i Ldurh) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldurh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}

func decodeLdurhOf(enc uint32, kind memKind, fp string) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		var rt string
		switch fp {
		case "s":
			rt = fpRegNameS(w & 0x1f)
		case "d":
			rt = fpRegNameD(w & 0x1f)
		case "x":
			rt = regNameX(w & 0x1f)
		case "w":
			rt = regNameW(w & 0x1f)
		default:
			rt = armRegName(w&0x1f, w>>30&3 == 3)
		}

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

		return Ldurh{
			base:   newBase(addr, w),
			lsBase: newLsBase(rt, rn, kind, off, tgt, enc, rm, option, shiftAmt),
		}
	}
}
