package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cmtst — cmtst.Arr vd, vn, vm (the three-same group: the lane width
// rides the arrangement).
type Cmtst struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b/4h/8h/2s/4s/2d
}

// newCmtst - the Cmtst constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newCmtst(b base, rd, rn, rm VReg, arr string) (Cmtst, error) {
	err := requireArr("Cmtst", arr, "8b", "16b", "4h", "8h", "2s", "4s", "2d")
	if err != nil {
		return Cmtst{}, err
	}

	return Cmtst{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const cmtstEnc uint32 = 237014016 // cmtst vd, vn, vm (Q=0, size=0 form)

func (i Cmtst) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cmtst.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Cmtst) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("cmtst: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("cmtst: %w", err)
	}

	return writeWord(w, cmtstEnc|q<<30|size<<22|rd|rn<<5|rm<<16)
}

func (Builder) Cmtst(rd, rn, rm VReg, arr string) (Instr, error) {
	return newCmtst(base{}, rd, rn, rm, arr)
}

func decodeCmtst(w uint32) (Instr, error) {
	in, err := newCmtst(
		newBase(w),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
