package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fdiv — fdiv fd, fn, fm (double/single by the operand kind).
type Fdiv struct {
	base

	rd, rn, rm string
}

// newFdiv - the Fdiv constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFdiv(b base, rd, rn, rm FReg) (Fdiv, error) {
	err := requireFpKind("Fdiv", rd, rn, rm)
	if err != nil {
		return Fdiv{}, err
	}

	return Fdiv{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const (
	fdivD uint32 = 0x1E601800 // fdiv dd, dn, dm
	fdivS uint32 = 0x1E201800 // fdiv sd, sn, sm
)

func (i Fdiv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fdiv %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Fdiv) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fdivD, fdivS)
	if err != nil {
		return 0, fmt.Errorf("fdiv: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fdiv: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Fdiv(rd, rn, rm FReg) (Instr, error) {
	return newFdiv(base{}, rd, rn, rm)
}

func decodeFdiv(w uint32) (Instr, error) {
	in, err := newFdiv(
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
