package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fneg — fneg fd, fn (double/single by the operand kind).
type Fneg struct {
	base

	rd, rn string
}

// newFneg - the Fneg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFneg(b base, rd, rn FReg) (Fneg, error) {
	err := requireFpKind("Fneg", rd, rn)
	if err != nil {
		return Fneg{}, err
	}

	return Fneg{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	fnegD uint32 = 0x1E614000 // fneg dd, dn
	fnegS uint32 = 0x1E214000 // fneg sd, sn
)

func (i Fneg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fneg %s, %s", i.rd, i.rn)
}

func (i Fneg) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fnegD, fnegS)
	if err != nil {
		return 0, fmt.Errorf("fneg: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fneg: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Fneg(rd, rn FReg) (Instr, error) {
	return newFneg(base{}, rd, rn)
}

func decodeFneg(w uint32) (Instr, error) {
	in, err := newFneg(
		newBase(w),
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
