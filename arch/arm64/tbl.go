package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tbl — tbl.16b vd, { vn }, vm.
type Tbl struct {
	base

	rd, rn, rm string
}

// newTbl - the Tbl constructor: the struct is assembled only
// here (NewTbl and the decoder call it).
func newTbl(b base, rd, rn, rm string) (Tbl, error) {
	return Tbl{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
	}, nil
}

// NewTbl - tbl.16b vd, { vn }, vm (the assembler ctor form; the word
// base is filled at decode time).
func NewTbl(rd, rn, rm string) (Tbl, error) {
	return newTbl(base{}, rd, rn, rm)
}

const tblEnc uint32 = 0x0E000000

func (i Tbl) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("tbl.16b %s, { %s }, %s", i.rd, i.rn, i.rm)
}

func (i Tbl) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("tbl: %w", err)
	}

	return writeWord(w, tblEnc|rd|rn<<5|rm<<16)
}

func decodeTbl(w uint32) (Instr, error) {
	in, err := newTbl(newBase(w), vReg(w&0x1f), vReg(w>>5&0x1f), vReg(w>>16&0x1f))
	if err != nil {
		return nil, err
	}

	return in, nil
}
