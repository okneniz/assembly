package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cmge — cmge.Arr vd, vn, vm (the three-same group: the lane width
// rides the arrangement).
type Cmge struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b/4h/8h/2s/4s/2d
}

// newCmge - the Cmge constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newCmge(b base, rd, rn, rm VReg, arr string) (Cmge, error) {
	err := requireArr("Cmge", arr, "8b", "16b", "4h", "8h", "2s", "4s", "2d")
	if err != nil {
		return Cmge{}, err
	}

	return Cmge{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const cmgeEnc uint32 = 236993536 // cmge vd, vn, vm (Q=0, size=0 form)

func (i Cmge) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cmge.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Cmge) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("cmge: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("cmge: %w", err)
	}

	return writeWord(w, cmgeEnc|q<<30|size<<22|rd|rn<<5|rm<<16)
}

func (Builder) Cmge(rd, rn, rm VReg, arr string) (Instr, error) {
	return newCmge(base{}, rd, rn, rm, arr)
}

func decodeCmge(w uint32) (Instr, error) {
	in, err := newCmge(
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
