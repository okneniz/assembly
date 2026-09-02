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

// Encodings of the unsigned-offset form: the access size is set by rt, the
// offset scale = log2 of the size.
const (
	ldrXEnc uint32 = 0xF9400000 // ldr xt, [xn, #imm12<<3]
	ldrWEnc uint32 = 0xB9400000 // ldr wt, [xn, #imm12<<2]
)

// Ldr — ldr rt, [rn, #off]: byte offset, scaling to the access size is
// hidden here.
func (Builder) Ldr(rt, rn Reg, off Off) (Instr, error) {
	if err := lsOperand(rt, rn, "Ldr"); err != nil {
		return nil, err
	}

	enc, scale := ldrXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = ldrWEnc, 2
	}

	if err := requireOff("Ldr", off, scale); err != nil {
		return nil, err
	}

	return Ldr{
		lsBase: newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	}, nil
}

func (i Ldr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldr %s, %s", i.rt, i.lsText())
}

func (i Ldr) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.lsWrite(w, pc, "ldr")
}

// ldrPoolWrap — pool ldr without self-verify: the decoder prints the slot's
// absolute address, not "=literal" — the text is not reproducible; the
// encoding is unambiguous (imm19 from the slot's known address).
type ldrPoolWrap struct {
	Ldr
}

func (ldrPoolWrap) SkipVerify() {}

func (i Ldr) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}

func decodeLdrOf(enc uint32, kind memKind, fp string) func(uint32, uint64) Instr {
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

		return Ldr{
			base:   newBase(addr, w),
			lsBase: newLsBase(rt, rn, kind, off, tgt, enc, rm, option, shiftAmt),
		}
	}
}
