package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ushr — ushr.Arr vd, vn, #shift (the shift lives in immh:immb,
// the stored value is esize-shift (printed shifted right)).
type Ushr struct {
	base

	rd, rn     string
	immh, immb uint32
	q          uint32
}

// newUshr - the Ushr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUshr(b base, q, immh, immb uint32, rd, rn VReg) (Ushr, error) {
	if immh > 0xf || immb > 7 {
		return Ushr{}, fmt.Errorf("arm64.NewUshr: imm out of range")
	}

	return Ushr{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immh: immh,
		immb: immb,
		q:    q,
	}, nil
}

const ushrEnc uint32 = 0x2F000400 // ushr vd, vn, #shift (Q=0 form)

func (i Ushr) ObjDump(_ disasm.ViewCtx) string {
	size, shift := decodeShiftImm(i.immh, i.immb)
	shift = uint32(8<<size) - shift
	return fmt.Sprintf("ushr.%s %s, %s, #0x%x",
		decodeArrangement(i.q, size), i.rd, i.rn, shift)
}

func (i Ushr) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("ushr: %w", err)
	}

	return writeWord(w, ushrEnc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}

// Ushr - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Ushr(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUshr: %w", err)
	}

	if shift == 0 || shift > 8<<size {
		return nil, fmt.Errorf("arm64.NewUshr: shift %d out of range for .%s",
			shift, arr)
	}

	// the stored value is esize-shift, under the size marker bit
	imm := 1<<(3+size) | (8<<size - shift)
	immh, immb := imm>>3, imm&7

	return newUshr(base{}, q, immh, immb, rd, rn)
}

func decodeUshr(w uint32) (Instr, error) {
	in, err := newUshr(
		newBase(w),
		w>>30&1,
		w>>19&0xf,
		w>>16&7,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
