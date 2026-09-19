package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Scvtf — scvtf fd, wn|xn (signed integer to FP; all width
// combinations are legal).
type Scvtf struct {
	base

	rd, rn string
}

// newScvtf - the Scvtf constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newScvtf(b base, rd FReg, rn Reg) (Scvtf, error) {
	err := requireGprClass(rn, "Scvtf", "rn")
	if err != nil {
		return Scvtf{}, err
	}

	return Scvtf{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	scvtfDX uint32 = 0x9E620000 // scvtf dd, xn
	scvtfDW uint32 = 0x1E620000 // scvtf dd, wn
	scvtfSX uint32 = 0x9E220000 // scvtf sd, xn
	scvtfSW uint32 = 0x1E220000 // scvtf sd, wn
)

func (i Scvtf) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("scvtf %s, %s", i.rd, i.rn)
}

func (i Scvtf) Encode(w io.Writer) (int64, error) {
	match := fpGprMatch(i.rd, i.rn, scvtfDX, scvtfDW, scvtfSX, scvtfSW)
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("scvtf: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Scvtf(rd FReg, rn Reg) (Instr, error) {
	return newScvtf(base{}, rd, rn)
}

func decodeScvtf(w uint32) (Instr, error) {
	in, err := newScvtf(
		newBase(w),
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// Ucvtf — ucvtf fd, wn|xn (unsigned integer to FP).
type Ucvtf struct {
	base

	rd, rn string
}

// newUcvtf - the Ucvtf constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newUcvtf(b base, rd FReg, rn Reg) (Ucvtf, error) {
	err := requireGprClass(rn, "Ucvtf", "rn")
	if err != nil {
		return Ucvtf{}, err
	}

	return Ucvtf{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	ucvtfDX uint32 = 0x9E630000 // ucvtf dd, xn
	ucvtfDW uint32 = 0x1E630000 // ucvtf dd, wn
	ucvtfSX uint32 = 0x9E230000 // ucvtf sd, xn
	ucvtfSW uint32 = 0x1E230000 // ucvtf sd, wn
)

func (i Ucvtf) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("ucvtf %s, %s", i.rd, i.rn)
}

func (i Ucvtf) Encode(w io.Writer) (int64, error) {
	match := fpGprMatch(i.rd, i.rn, ucvtfDX, ucvtfDW, ucvtfSX, ucvtfSW)
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("ucvtf: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Ucvtf(rd FReg, rn Reg) (Instr, error) {
	return newUcvtf(base{}, rd, rn)
}

func decodeUcvtf(w uint32) (Instr, error) {
	in, err := newUcvtf(
		newBase(w),
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
