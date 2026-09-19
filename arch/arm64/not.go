package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Not — not.8b/16b vd, vn (bitwise NOT per byte lane).
type Not struct {
	base

	rd, rn string
	arr    string // 8b/16b
}

// newNot - the Not constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newNot(b base, rd, rn VReg, arr string) (Not, error) {
	err := requireArr("Not", arr, "8b", "16b")
	if err != nil {
		return Not{}, err
	}

	return Not{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		arr:  arr,
	}, nil
}

const notEnc uint32 = 0x2E205800 // not vd, vn (Q=0 form)

func (i Not) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("not.%s %s, %s", i.arr, i.rd, i.rn)
}

func (i Not) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("not: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("not: %w", err)
	}

	return writeWord(w, notEnc|q<<30|size<<22|rd|rn<<5)
}

func (Builder) Not(rd, rn VReg, arr string) (Instr, error) {
	return newNot(base{}, rd, rn, arr)
}

func decodeNot(w uint32) (Instr, error) {
	in, err := newNot(
		newBase(w),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		decodeArrangement(w>>30&1, w>>22&3),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
