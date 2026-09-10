package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bfm — bfm rd, rn, #immr, #imms (no aliases).
type Bfm struct {
	base

	rd, rn     string
	immr, imms uint32
	isf        bool
}

// newBfm - the Bfm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newBfm(b base, rd Reg, rn Reg, immr uint32, imms uint32) (Bfm, error) {
	err := requireClass(
		rd,
		"Bfm",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Bfm{}, err
	}

	err = requireClass(
		rn,
		"Bfm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Bfm{}, err
	}

	err = requireWidth(
		"Bfm",
		rd,
		rn,
	)

	if err != nil {
		return Bfm{}, err
	}

	if immr > 63 || imms > 63 {
		return Bfm{}, fmt.Errorf(
			"arm64.NewBfm: operands immr/imms: %d/%d are out of 0..63",
			immr,
			imms,
		)
	}

	return Bfm{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immr: immr,
		imms: imms,
		isf:  rd.Is64(),
	}, nil
}

const (
	bfmX uint32 = 0xB3400000
	bfmW uint32 = 0x33000000
)

func (i Bfm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bfm %s, %s, #%d, #%d", i.rd, i.rn, i.immr, i.imms)
}

func (i Bfm) Encode(w io.Writer) (int64, error) {
	return bfmWrite(w, bfmX, bfmW, i.isf, i.rd, i.rn, i.immr, i.imms)
}

func (Builder) Bfm(rd, rn Reg, immr, imms uint32) (Instr, error) {
	return newBfm(base{}, rd, rn, immr, imms)
}

func decodeBfmInstr(w uint32) (Instr, error) {
	in, err := newBfm(
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
