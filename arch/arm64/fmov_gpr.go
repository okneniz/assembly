package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmovFromGpr — fmov fd, xn | fmov sn, wn (register bits moved from
// the integer file; xzr/wzr zeroes the FP register).
type FmovFromGpr struct {
	rd, rn string
}

// newFmovFromGpr - the FmovFromGpr constructor: validates the
// operands and assembles the struct (the Builder method delegates
// here; the decoder calls it with values read from the word).
func newFmovFromGpr(rd FReg, rn Reg) (FmovFromGpr, error) {
	err := requireGprClass(rn, "FmovFromGpr", "rn")
	if err != nil {
		return FmovFromGpr{}, err
	}

	err = requireFpGprWidth("FmovFromGpr", rd, rn)
	if err != nil {
		return FmovFromGpr{}, err
	}

	return FmovFromGpr{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

const (
	fmovFromGprD uint32 = 0x9E670000 // fmov dd, xn
	fmovFromGprS uint32 = 0x1E270000 // fmov sn, wn
)

func (i FmovFromGpr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmov %s, %s", i.rd, i.rn)
}

func (i FmovFromGpr) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fmovFromGprD, fmovFromGprS)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

// FmovToGpr — fmov xn, fd | fmov wn, sn (register bits moved to the
// integer file).
type FmovToGpr struct {
	rd, rn string
}

// newFmovToGpr - the FmovToGpr constructor: validates the operands
// and assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmovToGpr(rd Reg, rn FReg) (FmovToGpr, error) {
	err := requireGprClass(rd, "FmovToGpr", "rd")
	if err != nil {
		return FmovToGpr{}, err
	}

	err = requireFpGprWidth("FmovToGpr", rn, rd)
	if err != nil {
		return FmovToGpr{}, err
	}

	return FmovToGpr{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

const (
	fmovToGprX uint32 = 0x9E660000 // fmov xn, dd
	fmovToGprW uint32 = 0x1E260000 // fmov wn, sn
)

func (i FmovToGpr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmov %s, %s", i.rd, i.rn)
}

func (i FmovToGpr) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, fmovToGprX, fmovToGprW)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}
