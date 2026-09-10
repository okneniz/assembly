package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sturh — sturh ... (see lsBase for addressing kinds).
type Sturh struct {
	base
	lsBase
}

// newSturhBase - the Sturh constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newSturhBase(b base, e lsBase) Sturh {
	return Sturh{
		base:   b,
		lsBase: e,
	}
}

// newSturh - the Sturh constructor: validates the operands,
// delegates the assembly to newSturhBase.
func newSturh(b base, rt, rn Reg, off Off) (Sturh, error) {
	err := requireClass(
		rt,
		"Sturh",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Sturh{}, err
	}

	err = requireClass(
		rn,
		"Sturh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Sturh{}, err
	}

	if err = requireUnscaledOff("Sturh", off); err != nil {
		return Sturh{}, err
	}

	return newSturhBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, sturhEnc, "", "", 0),
	), nil
}

const sturhEnc uint32 = 0x78000000 // sturh wt, [xn, #±imm9]

func (i Sturh) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("sturh %s, %s", i.rt, i.lsText(ctx))
}

func (i Sturh) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "sturh")
}

func decodeSturhOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
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
		var lit int64
		var rm, option string
		var shiftAmt uint32
		switch kind {
		case memImm:
			off = int64(w>>10&0xfff) << (w >> 30 & 3)
		case memLiteral:
			lit = signExtendN(w>>5&0x7ffff, 19) * 4
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

		return Sturh{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Sturh(rt, rn Reg, off Off) (Instr, error) {
	return newSturh(base{}, rt, rn, off)
}
