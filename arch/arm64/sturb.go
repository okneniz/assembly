package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sturb — sturb ... (see lsBase for addressing kinds).
type Sturb struct {
	base
	lsBase
}

// newSturbBase - the Sturb constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newSturbBase(b base, e lsBase) Sturb {
	return Sturb{
		base:   b,
		lsBase: e,
	}
}

// newSturb - the Sturb constructor: validates the operands,
// delegates the assembly to newSturbBase.
func newSturb(b base, rt, rn Reg, off Off) (Sturb, error) {
	err := requireClass(
		rt,
		"Sturb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Sturb{}, err
	}

	err = requireClass(
		rn,
		"Sturb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Sturb{}, err
	}

	if err = requireUnscaledOff("Sturb", off); err != nil {
		return Sturb{}, err
	}

	return newSturbBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, sturbEnc, "", "", 0),
	), nil
}

const sturbEnc uint32 = 0x38000000 // sturb wt, [xn, #±imm9]

func (i Sturb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("sturb %s, %s", i.rt, i.lsText(ctx))
}

func (i Sturb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "sturb")
}

func decodeSturbOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
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

		return Sturb{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Sturb(rt, rn Reg, off Off) (Instr, error) {
	return newSturb(base{}, rt, rn, off)
}
