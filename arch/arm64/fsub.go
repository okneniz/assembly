package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fsub — fsub fd, fn, fm (double/single by the operand kind).
type Fsub struct {
	base

	rd, rn, rm string
}

// newFsub - the Fsub constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFsub(b base, rd, rn, rm FReg) (Fsub, error) {
	err := requireFpKind("Fsub", rd, rn, rm)
	if err != nil {
		return Fsub{}, err
	}

	return Fsub{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const (
	fsubD uint32 = 0x1E603800 // fsub dd, dn, dm
	fsubS uint32 = 0x1E203800 // fsub sd, sn, sm
)

func (i Fsub) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fsub %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Fsub) Encode(w io.Writer) (int64, error) {
	match, err := fpMatch(i.rd, fsubD, fsubS)
	if err != nil {
		return 0, fmt.Errorf("fsub: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fsub: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (Builder) Fsub(rd, rn, rm FReg) (Instr, error) {
	return newFsub(base{}, rd, rn, rm)
}

func decodeFsub(w uint32) (Instr, error) {
	in, err := newFsub(
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
