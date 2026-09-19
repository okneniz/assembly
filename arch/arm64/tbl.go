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
// here (the Builder method and the decoder call it).
func newTbl(b base, rd, rn, rm VReg) (Tbl, error) {
	return Tbl{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

func (Builder) Tbl(rd, rn, rm VReg) (Instr, error) {
	return newTbl(base{}, rd, rn, rm)
}

const tblEnc uint32 = 0x4E000000

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
	in, err := newTbl(newBase(w), newVReg(uint8(w&0x1f)), newVReg(uint8(w>>5&0x1f)), newVReg(uint8(w>>16&0x1f)))
	if err != nil {
		return nil, err
	}

	return in, nil
}
