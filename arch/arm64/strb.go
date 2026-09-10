package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Strb — strb ... (see lsBase for addressing kinds).
type Strb struct {
	base
	lsBase
}

// newStrbBase - the Strb constructor for a ready embedded base (the
// decoder and the string-operand layer): the struct is assembled only
// here.
func newStrbBase(b base, e lsBase) Strb {
	return Strb{
		base:   b,
		lsBase: e,
	}
}

// newStrb - the Strb constructor: validates the operands,
// delegates the assembly to newStrbBase.
func newStrb(b base, rt, rn Reg, off Off) (Strb, error) {
	err := requireClass(
		rt,
		"Strb",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Strb{}, err
	}

	err = requireClass(
		rn,
		"Strb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Strb{}, err
	}

	if err = requireOff("Strb", off, 0); err != nil {
		return Strb{}, err
	}

	return newStrbBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, strbEnc, "", "", 0),
	), nil
}

const strbEnc uint32 = 0x39000000 // strb wt, [xn, #imm12]

func (i Strb) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("strb %s, %s", i.rt, i.lsText(ctx))
}

func (i Strb) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "strb")
}

func decodeStrbOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameW(w & 0x1f)
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

		return Strb{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

// Strb — strb rt, [rn, #off]: byte access, rt — w register only
// (register 31 reads as wzr), rn — x register or SP (register 31 in the
// base reads as sp); the offset is an unscaled imm12 (0..0xfff).
func (Builder) Strb(rt, rn Reg, off Off) (Instr, error) {
	return newStrb(base{}, rt, rn, off)
}
