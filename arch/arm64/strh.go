package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Strh — strh ... (see lsBase for addressing kinds).
type Strh struct {
	base
	lsBase
}

// newStrhBase - the Strh constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newStrhBase(b base, e lsBase) Strh {
	return Strh{
		base:   b,
		lsBase: e,
	}
}

// newStrh - the Strh constructor: validates the operands,
// delegates the assembly to newStrhBase.
func newStrh(b base, rt, rn Reg, off Off) (Strh, error) {
	err := requireClass(
		rt,
		"Strh",
		"rt",
		"w register (register 31 in rt reads as wzr)",
		classW,
		classWZR,
	)

	if err != nil {
		return Strh{}, err
	}

	err = requireClass(
		rn,
		"Strh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Strh{}, err
	}

	if err = requireOff("Strh", off, 1); err != nil {
		return Strh{}, err
	}

	return newStrhBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, strhEnc, "", "", 0),
	), nil
}

const strhEnc uint32 = 0x79000000 // strh wt, [xn, #imm12<<1]

func (i Strh) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("strh %s, %s", i.rt, i.lsText(ctx))
}

func (i Strh) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "strh")
}

func decodeStrhOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
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

		return Strh{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Strh(rt, rn Reg, off Off) (Instr, error) {
	return newStrh(base{}, rt, rn, off)
}
