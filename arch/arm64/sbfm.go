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

// newSbfm - the Sbfm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newSbfm(b base, rd Reg, rn Reg, immr uint32, imms uint32) (Sbfm, error) {
	err := requireClass(
		rd,
		"Sbfm",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Sbfm{}, err
	}

	err = requireClass(
		rn,
		"Sbfm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Sbfm{}, err
	}

	err = requireWidth(
		"Sbfm",
		rd,
		rn,
	)

	if err != nil {
		return Sbfm{}, err
	}

	if immr > 63 || imms > 63 {
		return Sbfm{}, fmt.Errorf(
			"arm64.NewSbfm: operands immr/imms: %d/%d are out of 0..63",
			immr,
			imms,
		)
	}

	return Sbfm{
		base:  b,
		rd:    rd.name(),
		rn:    rn.name(),
		immr:  immr,
		imms:  imms,
		rnNum: rn.bits(),
		isf:   rd.Is64(),
	}, nil
}

const (
	sbfmX uint32 = 0x93400000
	sbfmW uint32 = 0x13000000
)

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

func (i Sbfm) Encode(w io.Writer) (int64, error) {
	return bfmWrite(w, sbfmX, sbfmW, i.isf, i.rd, i.rn, i.immr, i.imms)
}

func (Builder) Sbfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	return newSbfm(base{}, rd, rn, immr, imms)
}

func decodeSbfm(w uint32) (Instr, error) {
	in, err := newSbfm(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		w>>16&0x3f,
		w>>10&0x3f,
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
