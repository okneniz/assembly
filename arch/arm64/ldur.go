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

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldurXEnc uint32 = 0xF8400000 // ldur xt, [xn, #±imm9]
	ldurWEnc uint32 = 0xB8400000 // ldur wt, [xn, #±imm9]
)

// Ldur — ldur rt, [rn, #off]: the unscaled form, rt — x/w register
// (register 31 reads as zr), rn — x register or SP (register 31 in the
// base reads as sp); the offset is a signed imm9 (-0x100..0xff, any
// alignment).
func (Builder) Ldur(rt, rn Reg, off Off) (Instr, error) {
	if err := lsOperand(rt, rn, "Ldur"); err != nil {
		return nil, err
	}

	enc := ldurWEnc
	if rt.Is64() {
		enc = ldurXEnc
	}

	if err := requireUnscaledOff("Ldur", off); err != nil {
		return nil, err
	}

	return Ldur{
		lsBase: newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, enc, "", "", 0),
	}, nil
}

func (i Ldur) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldur %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldur) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldur")
}

func decodeLdurOf(enc uint32, kind memKind, fp string) func(uint32) Instr {
	return func(w uint32) Instr {
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
		}
	}
}
