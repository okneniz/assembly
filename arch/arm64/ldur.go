package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldur — ldur ... (see lsBase for the addressing kinds).
type Ldur struct {
	base
	lsBase
}

// newLdurBase - the Ldur constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdurBase(b base, e lsBase) Ldur {
	return Ldur{
		base:   b,
		lsBase: e,
	}
}

// newLdur - the Ldur constructor: validates the operands,
// delegates the assembly to newLdurBase.
func newLdur(b base, rt, rn Reg, off Off) (Ldur, error) {
	if err := lsOperand(rt, rn, "Ldur"); err != nil {
		return Ldur{}, err
	}

	enc := ldurWEnc
	if rt.Is64() {
		enc = ldurXEnc
	}

	if err := requireUnscaledOff("Ldur", off); err != nil {
		return Ldur{}, err
	}

	return newLdurBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldurXEnc uint32 = 0xF8400000 // ldur xt, [xn, #±imm9]
	ldurWEnc uint32 = 0xB8400000 // ldur wt, [xn, #±imm9]
)

func (i Ldur) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldur %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldur) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldur")
}

func decodeLdurOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
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

		return Ldur{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Ldur(rt, rn Reg, off Off) (Instr, error) {
	return newLdur(base{}, rt, rn, off)
}
