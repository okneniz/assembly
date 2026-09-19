package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fmov — fmov fd, fn (the register form between two FP registers of
// one kind; the GPR and immediate forms are FmovFromGpr/FmovToGpr/
// FmovImm).
type Fmov struct {
	base

	rd, rn string
}

// newFmov - the Fmov constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmov(b base, rd, rn FReg) (Fmov, error) {
	err := requireFpKind("Fmov", rd, rn)
	if err != nil {
		return Fmov{}, err
	}

	return Fmov{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	fmovD uint32 = 0x1E604000 // fmov dd, dn
	fmovS uint32 = 0x1E204000 // fmov sd, sn
)

func (i Fmov) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmov %s, %s", i.rd, i.rn)
}

func (i Fmov) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fmovD, fmovS)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Fmov(rd, rn FReg) (Instr, error) {
	return newFmov(base{}, rd, rn)
}

func decodeFmov(w uint32) (Instr, error) {
	in, err := newFmov(
		newBase(w),
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
