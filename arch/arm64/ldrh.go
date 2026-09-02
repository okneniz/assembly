package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldrh — ldrh ... (see lsBase for the addressing kinds).
type Ldrh struct {
	base
	lsBase
}

const ldrhEnc uint32 = 0x79400000 // ldrh wt, [xn, #imm12<<1]

// Ldrh — ldrh rt, [rn, #off]: halfword access, rt — w register
// only (register 31 reads as wzr), rn — x register or SP (register 31
// in the base reads as sp); the offset is an imm12 scaled by 2
// (0..0x1ffe, alignment 2).
func (Builder) Ldrh(rt, rn Reg, off Off) (Instr, error) {
	if err := requireClass(rt, "Ldrh", "rt", "w register (register 31 in rt reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Ldrh",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	if err := requireOff("Ldrh", off, 1); err != nil {
		return nil, err
	}

	return Ldrh{
		lsBase: newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, ldrhEnc, "", "", 0),
	}, nil
}

func (i Ldrh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ldrh %s, %s", i.rt, i.lsText())
}

func (i Ldrh) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.lsWrite(w, pc, "ldrh")
}

func (i Ldrh) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldrh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}

func decodeLdrhOf(enc uint32, kind memKind) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		rt := regNameW(w & 0x1f)
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

		return Ldrh{
			base:   newBase(addr, w),
			lsBase: newLsBase(rt, rn, kind, off, tgt, enc, rm, option, shiftAmt),
		}
	}
}
