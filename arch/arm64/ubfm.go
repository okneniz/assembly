package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ubfm — ubfm rd, rn, #immr, #imms; aliases lsl/lsr/ubfiz/ubfx.
type Ubfm struct {
	base

	rd, rn     string
	immr, imms uint32
	isf        bool
}

// newUbfm - the Ubfm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newUbfm(b base, rd Reg, rn Reg, immr uint32, imms uint32) (Ubfm, error) {
	err := requireClass(
		rd,
		"Ubfm",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Ubfm{}, err
	}

	err = requireClass(
		rn,
		"Ubfm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Ubfm{}, err
	}

	err = requireWidth(
		"Ubfm",
		rd,
		rn,
	)

	if err != nil {
		return Ubfm{}, err
	}

	if immr > 63 || imms > 63 {
		return Ubfm{}, fmt.Errorf(
			"arm64.NewUbfm: operands immr/imms: %d/%d are out of 0..63",
			immr,
			imms,
		)
	}

	return Ubfm{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immr: immr,
		imms: imms,
		isf:  rd.Is64(),
	}, nil
}

const (
	ubfmX uint32 = 0xD3400000
	ubfmW uint32 = 0x53000000
)

func (i Ubfm) ObjDump(_ disasm.ViewCtx) string {
	regsize := bfmRegsize(i.rd, i.immr, i.imms)
	if i.imms != regsize-1 && i.imms+1 == i.immr { // LSL alias
		return fmt.Sprintf("lsl %s, %s, #%d", i.rd, i.rn, regsize-i.immr)
	}

	if i.imms == regsize-1 { // LSR alias
		return fmt.Sprintf("lsr %s, %s, #%d", i.rd, i.rn, i.immr)
	}

	if i.imms < i.immr { // UBFIZ: lsb = regsize-immr, width = imms+1
		return fmt.Sprintf("ubfiz %s, %s, #%d, #%d", i.rd, i.rn, regsize-i.immr, i.imms+1)
	}

	return fmt.Sprintf("ubfx %s, %s, #%d, #%d", i.rd, i.rn, i.immr, i.imms-i.immr+1)
}

func (i Ubfm) Encode(w io.Writer) (int64, error) {
	return bfmWrite(w, ubfmX, ubfmW, i.isf, i.rd, i.rn, i.immr, i.imms)
}

func (Builder) Ubfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	return newUbfm(base{}, rd, rn, immr, imms)
}

func decodeUbfm(w uint32) (Instr, error) {
	in, err := newUbfm(
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
