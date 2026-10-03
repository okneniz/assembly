package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Extr — extr rd, rn, rm, #lsb; pseudo: ror rd, rn, #imm (rn == rm).
type Extr struct {
	rd, rn, rm string
	lsb        uint32
	isf        bool
}

// newExtr - the Extr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newExtr(rd Reg, rn Reg, rm Reg, lsb Imm6) (Extr, error) {
	for _, r := range []struct {
		reg Reg
		op  string
	}{{
		rd,
		"rd",
	}, {
		rn,
		"rn",
	}, {
		rm,
		"rm",
	}} {
		err := requireClass(
			r.reg,
			"Extr",
			r.op,
			"register 31 reads as zr — use XZR/WZR",
			classX,
			classW,
			classXZR,
			classWZR,
		)

		if err != nil {
			return Extr{}, err
		}
	}

	if err := requireWidth("Extr", rd, rn, rm); err != nil {
		return Extr{}, err
	}

	isf := rd.Is64()
	if !isf && lsb.v > 31 {
		return Extr{}, errors.New("arm64.NewExtr: lsb out of range for the 32-bit form")
	}

	return Extr{
		rd:  rd.name(),
		rn:  rn.name(),
		rm:  rm.name(),
		lsb: lsb.v,
		isf: isf,
	}, nil
}

const (
	extrW uint32 = 0x13800000 // sf=0, N=0
	extrX uint32 = 0x93C00000 // sf=1, N=1 (N must equal sf)
)

func (i Extr) Encode(w io.Writer) (int64, error) {
	match := extrX
	if !i.isf {
		match = extrW
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("extr: %w", err)
	}

	if i.lsb > 63 || (!i.isf && i.lsb > 31) {
		return 0, fmt.Errorf("extr: lsb %#x out of range", i.lsb)
	}

	return writeWord(w, match|rd|rn<<5|i.lsb<<10|rm<<16)
}

func (i Extr) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == i.rm {
		return fmt.Sprintf("ror %s, %s, #0x%x", i.rd, i.rn, i.lsb)
	}

	return fmt.Sprintf("extr %s, %s, %s, #0x%x", i.rd, i.rn, i.rm, i.lsb)
}
