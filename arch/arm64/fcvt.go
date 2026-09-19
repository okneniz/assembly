package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fcvt — fcvt fd, fn (the width conversion s↔d; the operand kinds
// must differ — the source type rides bits [22:21], the destination
// type the opcode bits [15:14]).
type Fcvt struct {
	base

	rd, rn string
}

// newFcvt - the Fcvt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFcvt(b base, rd, rn FReg) (Fcvt, error) {
	if rd.Is64() == rn.Is64() {
		return Fcvt{}, fmt.Errorf(
			"arm64.NewFcvt: converts between s and d - kinds must differ: %s vs %s",
			rd.name(), rn.name(),
		)
	}

	return Fcvt{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const (
	fcvtSD uint32 = 0x1E624000 // fcvt sd, dn
	fcvtDS uint32 = 0x1E22C000 // fcvt dd, sn
)

func (i Fcvt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fcvt %s, %s", i.rd, i.rn)
}

func (i Fcvt) Encode(w io.Writer) (int64, error) {
	match := fcvtSD
	if i.rd[0] == 'd' {
		match = fcvtDS
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fcvt: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (Builder) Fcvt(rd, rn FReg) (Instr, error) {
	return newFcvt(base{}, rd, rn)
}

func decodeFcvt(w uint32) (Instr, error) {
	in, err := newFcvt(
		newBase(w),
		newFReg(uint8(w&0x1f), w>>15&1 == 1),
		newFReg(uint8(w>>5&0x1f), w>>22&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
