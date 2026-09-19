package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fmax — fmax fd, fn, fm (double/single by the operand kind).
type Fmax struct {
	base

	rd, rn, rm string
}

// newFmax - the Fmax constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmax(b base, rd, rn, rm FReg) (Fmax, error) {
	err := requireFpKind("Fmax", rd, rn, rm)
	if err != nil {
		return Fmax{}, err
	}

	return Fmax{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const (
	fmaxD uint32 = 0x1E604800 // fmax dd, dn, dm
	fmaxS uint32 = 0x1E204800 // fmax sd, sn, sm
)

func (i Fmax) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmax %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Fmax) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fmaxD, fmaxS)
	if err != nil {
		return 0, fmt.Errorf("fmax: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmax: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Fmax(rd, rn, rm FReg) (Instr, error) {
	return newFmax(base{}, rd, rn, rm)
}

func decodeFmax(w uint32) (Instr, error) {
	in, err := newFmax(
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
