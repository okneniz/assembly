package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldurb — ldurb ... (see lsBase for the addressing kinds).
type Ldurb struct {
	base
	lsBase
}

// newLdurbBase - the Ldurb constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdurbBase(b base, e lsBase) Ldurb {
	return Ldurb{
		base:   b,
		lsBase: e,
	}
}

// newLdurb - the Ldurb constructor: validates the operands,
// delegates the assembly to newLdurbBase.
func newLdurb(b base, rt, rn Reg, off Off) (Ldurb, error) {
	err := requireClass(
		rt,
		"Ldurb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Ldurb{}, err
	}

	err = requireClass(
		rn,
		"Ldurb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldurb{}, err
	}

	if err = requireUnscaledOff("Ldurb", off); err != nil {
		return Ldurb{}, err
	}

	return newLdurbBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, ldurbEnc, "", "", 0),
	), nil
}

const ldurbEnc uint32 = 0x38400000 // ldurb wt, [xn, #±imm9]

func (i Ldurb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldurb %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldurb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldurb")
}

func decodeLdurbOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
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

		return Ldurb{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Ldurb(rt, rn Reg, off Off) (Instr, error) {
	return newLdurb(base{}, rt, rn, off)
}
