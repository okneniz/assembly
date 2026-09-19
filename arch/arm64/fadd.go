package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fadd — fadd fd, fn, fm (double/single by the operand kind).
type Fadd struct {
	base

	rd, rn, rm string
}

// newFadd - the Fadd constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFadd(b base, rd, rn, rm FReg) (Fadd, error) {
	err := requireFpKind("Fadd", rd, rn, rm)
	if err != nil {
		return Fadd{}, err
	}

	return Fadd{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const (
	faddD uint32 = 0x1E602800 // fadd dd, dn, dm
	faddS uint32 = 0x1E202800 // fadd sd, sn, sm
)

func (i Fadd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fadd %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Fadd) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, faddD, faddS)
	if err != nil {
		return 0, fmt.Errorf("fadd: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fadd: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Fadd(rd, rn, rm FReg) (Instr, error) {
	return newFadd(base{}, rd, rn, rm)
}

func decodeFadd(w uint32) (Instr, error) {
	in, err := newFadd(
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
