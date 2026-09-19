package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fmul — fmul fd, fn, fm (double/single by the operand kind).
type Fmul struct {
	base

	rd, rn, rm string
}

// newFmul - the Fmul constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmul(b base, rd, rn, rm FReg) (Fmul, error) {
	err := requireFpKind("Fmul", rd, rn, rm)
	if err != nil {
		return Fmul{}, err
	}

	return Fmul{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const (
	fmulD uint32 = 0x1E600800 // fmul dd, dn, dm
	fmulS uint32 = 0x1E200800 // fmul sd, sn, sm
)

func (i Fmul) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmul %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Fmul) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fmulD, fmulS)
	if err != nil {
		return 0, fmt.Errorf("fmul: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmul: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Fmul(rd, rn, rm FReg) (Instr, error) {
	return newFmul(base{}, rd, rn, rm)
}

func decodeFmul(w uint32) (Instr, error) {
	in, err := newFmul(
		newBase(w),
		newFReg(uint8(w&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
		newFReg(uint8(w>>16&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
