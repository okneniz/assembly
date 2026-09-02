package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sbfm — sbfm rd, rn, #immr, #imms; aliases asr/sxtb/sxth/sxtw/sbfiz/sbfx.
type Sbfm struct {
	base

	rd, rn     string
	immr, imms uint32
	rnNum      uint32 // for sxt*: the W name of Rn
	isf        bool
}

const (
	sbfmX uint32 = 0x93400000
	sbfmW uint32 = 0x13000000
)

// Sbfm — sbfm rd, rn, #immr, #imms (asr/sxtb/sxth/sxtw/sbfiz/sbfx
// aliases — the printed form depends on immr/imms). Register 31 reads as
// zr (SP/WSP are not allowed — use XZR/WZR); the width is shared by both
// registers; immr/imms — 0..63 (the top half of the range is
// unpredictable in the 32-bit form, as in the architecture).
func (Builder) Sbfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	if err := requireClass(rd, "Sbfm", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Sbfm", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("Sbfm", rd, rn); err != nil {
		return nil, err
	}

	if immr > 63 || imms > 63 {
		return nil, fmt.Errorf(
			"arm64.NewSbfm: operands immr/imms: %d/%d are out of 0..63",
			immr,
			imms,
		)
	}

	return Sbfm{
		rd:    rd.name(),
		rn:    rn.name(),
		immr:  immr,
		imms:  imms,
		rnNum: rn.bits(),
		isf:   rd.Is64(),
	}, nil
}

func decodeSbfm(w uint32, addr uint64) Instr {
	return Sbfm{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		rn:    armRegName(w>>5&0x1f, w>>31&1 == 1),
		immr:  w >> 16 & 0x3f,
		imms:  w >> 10 & 0x3f,
		rnNum: w >> 5 & 0x1f,
		isf:   w>>31&1 == 1,
	}
}

func (i Sbfm) ObjDump(_ disasm.ViewCtx) string {
	regsize := bfmRegsize(i.rd, i.immr, i.imms)
	// LSL is a UBFM alias, not SBFM: LLVM/GNU print SBFM encodings with
	// imms+1 == immr as sbfiz (the SBFIZ condition, imms < immr, covers them).
	if i.imms == regsize-1 { // ASR alias
		return fmt.Sprintf("asr %s, %s, #%d", i.rd, i.rn, i.immr)
	}

	// SXTB/SXTH/SXTW with immr=0 and imms=7/15/31 (regsize=64): Rn is a W register.
	if i.immr == 0 && regsize == 64 {
		rnW := regNameW(i.rnNum)
		switch i.imms {
		case 7:
			return fmt.Sprintf("sxtb %s, %s", i.rd, rnW)
		case 15:
			return fmt.Sprintf("sxth %s, %s", i.rd, rnW)
		case 31:
			return fmt.Sprintf("sxtw %s, %s", i.rd, rnW)
		}
	}

	if i.imms < i.immr { // SBFIZ
		return fmt.Sprintf("sbfiz %s, %s, #%d, #%d", i.rd, i.rn, regsize-i.immr, i.imms+1)
	}

	return fmt.Sprintf("sbfx %s, %s, #%d, #%d", i.rd, i.rn, i.immr, i.imms-i.immr+1)
}

func (i Sbfm) Encode(w io.Writer, pc uint64) (int64, error) {
	return bfmWrite(w, sbfmX, sbfmW, i.isf, i.rd, i.rn, i.immr, i.imms)
}

func (i Sbfm) MarshalJSON() ([]byte, error) {
	return i.marshal("sbfm", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}
