package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Madd — madd rd, rn, rm, ra; pseudo: mul (ra = xzr).
type Madd struct {
	base

	rd, rn, rm, ra string
}

// newMadd - the Madd constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newMadd(b base, rd Reg, rn Reg, rm Reg, ra Reg) (Madd, error) {
	err := requireClass(
		rd,
		"Madd",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Madd{}, err
	}

	err = requireClass(
		rn,
		"Madd",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Madd{}, err
	}

	err = requireClass(
		rm,
		"Madd",
		"rm",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Madd{}, err
	}

	err = requireClass(
		ra,
		"Madd",
		"ra",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Madd{}, err
	}

	err = requireWidth(
		"Madd",
		rd,
		rn,
		rm,
		ra,
	)

	if err != nil {
		return Madd{}, err
	}

	return Madd{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		ra:   ra.name(),
		rm:   rm.name(),
	}, nil
}

const (
	maddX uint32 = 0x9B000000
	maddW uint32 = 0x1B000000
)

func (i Madd) ObjDump(_ disasm.ViewCtx) string {
	zr := "xzr"
	if i.rd[0] == 'w' {
		zr = "wzr"
	}

	if i.ra == zr {
		return fmt.Sprintf("mul %s, %s, %s", i.rd, i.rn, i.rm)
	}

	return fmt.Sprintf("madd %s, %s, %s, %s", i.rd, i.rn, i.rm, i.ra)
}

func (i Madd) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, maddX, maddW)
	if err != nil {
		return 0, fmt.Errorf("madd: %w", err)
	}

	return msubWrite(w, match, i)
}

// msubWrite - the shared word of the madd/msub family (msub = Madd with the opcode bit).
func msubWrite(w io.Writer, match uint32, i Madd) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, err
	}

	ra, err := armRegNum(i.ra)
	if err != nil {
		return 0, err
	}

	return writeWord(w, match|rd|rn<<5|ra<<10|rm<<16)
}

func (Builder) Madd(rd, rn, rm, ra Reg) (Instr, error) {
	return newMadd(base{}, rd, rn, rm, ra)
}

func decodeMadd(w uint32) (Instr, error) {
	in, err := newMadd(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
		gprOf(w>>10&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
