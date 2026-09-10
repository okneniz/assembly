package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrsw — ldrsw ... (see lsBase for the addressing kinds).
type Ldrsw struct {
	base
	lsBase
}

// newLdrswBase - the Ldrsw constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdrswBase(b base, e lsBase) Ldrsw {
	return Ldrsw{
		base:   b,
		lsBase: e,
	}
}

// newLdrsw - the Ldrsw constructor: validates the operands,
// delegates the assembly to newLdrswBase.
func newLdrsw(b base, rt, rn Reg, off Off) (Ldrsw, error) {
	err := requireClass(
		rt,
		"Ldrsw",
		"rt",
		"x register (register 31 in rt reads as xzr)",
		classX,
		classXZR,
	)

	if err != nil {
		return Ldrsw{}, err
	}

	err = requireClass(
		rn,
		"Ldrsw",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Ldrsw{}, err
	}

	if err = requireOff("Ldrsw", off, 2); err != nil {
		return Ldrsw{}, err
	}

	return newLdrswBase(
		b,
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, ldrswEnc, "", "", 0),
	), nil
}

const ldrswEnc uint32 = 0xB9800000 // ldrsw xt, [xn, #imm12<<2]

func (i Ldrsw) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldrsw %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldrsw) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldrsw")
}

func decodeLdrswOf(enc uint32, kind memKind) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		rt := regNameX(w & 0x1f)
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

		return Ldrsw{
			base:   newBase(w),
			lsBase: newLsBase(rt, rn, kind, off, lit, enc, rm, option, shiftAmt),
		}, nil
	}
}

func (Builder) Ldrsw(rt, rn Reg, off Off) (Instr, error) {
	return newLdrsw(base{}, rt, rn, off)
}
