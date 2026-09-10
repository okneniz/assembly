package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Str — str ... (see lsBase for addressing kinds).
type Str struct {
	base
	lsBase
}

// newStrBase - the Str constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newStrBase(b base, e lsBase) Str {
	return Str{
		base:   b,
		lsBase: e,
	}
}

// newStr - the Str constructor: validates the operands,
// delegates the assembly to newStrBase.
func newStr(b base, rt, rn Reg, off Off) (Str, error) {
	if err := lsOperand(rt, rn, "Str"); err != nil {
		return Str{}, err
	}

	enc, scale := strXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = strWEnc, 2
	}

	if err := requireOff("Str", off, scale); err != nil {
		return Str{}, err
	}

	return newStrBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the unsigned-offset form: the access size is set by rt, the
// offset scale = log2 of the size.
const (
	strXEnc uint32 = 0xF9000000 // str xt, [xn, #imm12<<3]
	strWEnc uint32 = 0xB9000000 // str wt, [xn, #imm12<<2]
)

func (i Str) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("str %s, %s", i.rt, i.lsText(ctx))
}

func (i Str) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "str")
}

func decodeStrOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
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

		return Str{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Str(rt, rn Reg, off Off) (Instr, error) {
	return newStr(base{}, rt, rn, off)
}
