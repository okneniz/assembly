package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stur — stur ... (see lsBase for addressing kinds).
type Stur struct {
	base
	lsBase
}

// newSturBase - the Stur constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newSturBase(b base, e lsBase) Stur {
	return Stur{
		base:   b,
		lsBase: e,
	}
}

// newStur - the Stur constructor: validates the operands,
// delegates the assembly to newSturBase.
func newStur(b base, rt, rn Reg, off Off) (Stur, error) {
	if err := lsOperand(rt, rn, "Stur"); err != nil {
		return Stur{}, err
	}

	enc := sturWEnc
	if rt.Is64() {
		enc = sturXEnc
	}

	if err := requireUnscaledOff("Stur", off); err != nil {
		return Stur{}, err
	}

	return newSturBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	sturXEnc uint32 = 0xF8000000 // stur xt, [xn, #±imm9]
	sturWEnc uint32 = 0xB8000000 // stur wt, [xn, #±imm9]
)

func (i Stur) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("stur %s, %s", i.rt, i.lsText(ctx))
}

func (i Stur) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "stur")
}

func decodeSturOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
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

		return Stur{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Stur(rt, rn Reg, off Off) (Instr, error) {
	return newStur(base{}, rt, rn, off)
}
