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

const ldurbEnc uint32 = 0x38400000 // ldurb wt, [xn, #±imm9]

// Ldurb — ldurb rt, [rn, #off]: the unscaled form, byte access,
// rt — w register only (register 31 reads as wzr), rn — x register or
// SP (register 31 in the base reads as sp); the offset is a signed
// imm9 (-0x100..0xff, any alignment).
func (Builder) Ldurb(rt, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldurb", "rt", "w register (register 31 in rt reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldurb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireUnscaledOff("Ldurb", off); err != nil {
		return nil, err
	}

	return Ldurb{
		lsBase: newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, ldurbEnc, "", "", 0),
	}, nil
}

func (i Ldurb) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldurb %s, %s", i.rt, i.lsText())
}

func (i Ldurb) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.lsWrite(w, pc, "ldurb")
}

func (i Ldurb) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldurb",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}

func decodeLdurbOf(enc uint32, kind memKind, fp string) func(uint32, uint64) Instr {
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

		return Ldurb{
			base:   newBase(addr, w),
			lsBase: newLsBase(rt, rn, kind, off, tgt, enc, rm, option, shiftAmt),
		}
	}
}
