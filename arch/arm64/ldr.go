package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldr — ldr ... (see lsBase for the addressing kinds).
type Ldr struct {
	base
	lsBase
}

// newLdrBase - the Ldr constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdrBase(b base, e lsBase) Ldr {
	return Ldr{
		base:   b,
		lsBase: e,
	}
}

// newLdr - the Ldr constructor: validates the operands,
// delegates the assembly to newLdrBase.
func newLdr(b base, rt, rn Reg, off Off) (Ldr, error) {
	if err := lsOperand(rt, rn, "Ldr"); err != nil {
		return Ldr{}, err
	}

	enc, scale := ldrXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = ldrWEnc, 2
	}

	if err := requireOff("Ldr", off, scale); err != nil {
		return Ldr{}, err
	}

	return newLdrBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the unsigned-offset form: the access size is set by rt, the
// offset scale = log2 of the size.
const (
	ldrXEnc uint32 = 0xF9400000 // ldr xt, [xn, #imm12<<3]
	ldrWEnc uint32 = 0xB9400000 // ldr wt, [xn, #imm12<<2]
)

func (i Ldr) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldr %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldr) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldr")
}

// ldrPoolWrap — pool ldr without self-verify: the decoder prints the slot's
// absolute address, not "=literal" — the text is not reproducible; the
// encoding is unambiguous (imm19 from the slot's known address).
type ldrPoolWrap struct {
	Ldr
}

func (ldrPoolWrap) SkipVerify() {}

func decodeLdrOf(enc uint32, kind memKind, fp string) func(uint32) (Instr, error) {
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

		return Ldr{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

// LdrPoolWrapOf — the pool-wrapped literal ldr (PoolUser); lit — the
// pc-relative byte offset of the pool slot.
func LdrPoolWrapOf(rt string, lit int64, enc uint32) Instr {
	return ldrPoolWrap{Ldr{lsBase: newLsBase(rt, "", memLiteral, 0, lit, enc, "", "", 0)}}
}

func (Builder) Ldr(rt, rn Reg, off Off) (Instr, error) {
	return newLdr(base{}, rt, rn, off)
}
