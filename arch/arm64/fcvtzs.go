package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fcvtzs — fcvtzs wd|xd, fn (FP to signed integer; all width
// combinations are legal).
type Fcvtzs struct {
	base

	rd, rn string
}

// newFcvtzs - the Fcvtzs constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFcvtzs(b base, rd Reg, rn FReg) (Fcvtzs, error) {
	err := requireGprClass(rd, "Fcvtzs", "rd")
	if err != nil {
		return Fcvtzs{}, err
	}

	return Fcvtzs{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	fcvtzsXD uint32 = 0x9E780000 // fcvtzs xd, dn
	fcvtzsWD uint32 = 0x1E780000 // fcvtzs wd, dn
	fcvtzsXS uint32 = 0x9E380000 // fcvtzs xd, sn
	fcvtzsWS uint32 = 0x1E380000 // fcvtzs wd, sn
)

func (i Fcvtzs) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fcvtzs %s, %s", i.rd, i.rn)
}

func (i Fcvtzs) Encode(w io.Writer) (int64, error) {
	match := fpGprMatch(i.rn, i.rd, fcvtzsXD, fcvtzsWD, fcvtzsXS, fcvtzsWS)
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fcvtzs: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Fcvtzs(rd Reg, rn FReg) (Instr, error) {
	return newFcvtzs(base{}, rd, rn)
}

func decodeFcvtzs(w uint32) (Instr, error) {
	in, err := newFcvtzs(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// Fcvtzu — fcvtzu wd|xd, fn (FP to unsigned integer).
type Fcvtzu struct {
	base

	rd, rn string
}

// newFcvtzu - the Fcvtzu constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFcvtzu(b base, rd Reg, rn FReg) (Fcvtzu, error) {
	err := requireGprClass(rd, "Fcvtzu", "rd")
	if err != nil {
		return Fcvtzu{}, err
	}

	return Fcvtzu{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	fcvtzuXD uint32 = 0x9E790000 // fcvtzu xd, dn
	fcvtzuWD uint32 = 0x1E790000 // fcvtzu wd, dn
	fcvtzuXS uint32 = 0x9E390000 // fcvtzu xd, sn
	fcvtzuWS uint32 = 0x1E390000 // fcvtzu wd, sn
)

func (i Fcvtzu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fcvtzu %s, %s", i.rd, i.rn)
}

func (i Fcvtzu) Encode(w io.Writer) (int64, error) {
	match := fpGprMatch(i.rn, i.rd, fcvtzuXD, fcvtzuWD, fcvtzuXS, fcvtzuWS)
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fcvtzu: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Fcvtzu(rd Reg, rn FReg) (Instr, error) {
	return newFcvtzu(base{}, rd, rn)
}

func decodeFcvtzu(w uint32) (Instr, error) {
	in, err := newFcvtzu(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
