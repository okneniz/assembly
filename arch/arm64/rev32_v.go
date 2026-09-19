package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rev32V — rev32.8b/16b/4h/8h vd, vn (byte order reversed within 32-bit containers).
type Rev32V struct {
	base

	rd, rn string
	arr    string // 8b/16b/4h/8h
}

// newRev32V - the Rev32V constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newRev32V(b base, rd, rn VReg, arr string) (Rev32V, error) {
	err := requireArr("Rev32V", arr, "8b", "16b", "4h", "8h")
	if err != nil {
		return Rev32V{}, err
	}

	return Rev32V{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		arr:  arr,
	}, nil
}

const rev32VEnc uint32 = 0x2E200800 // rev32 vd, vn (Q=0 form)

func (i Rev32V) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev32.%s %s, %s", i.arr, i.rd, i.rn)
}

func (i Rev32V) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rev32: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("rev32: %w", err)
	}

	return writeWord(w, rev32VEnc|q<<30|size<<22|rd|rn<<5)
}

func (Builder) Rev32V(rd, rn VReg, arr string) (Instr, error) {
	return newRev32V(base{}, rd, rn, arr)
}

func decodeRev32V(w uint32) (Instr, error) {
	in, err := newRev32V(
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
