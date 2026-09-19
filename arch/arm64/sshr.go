package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sshr — sshr.Arr vd, vn, #shift (the shift lives in immh:immb,
// the stored value is esize-shift (printed shifted right)).
type Sshr struct {
	base

	rd, rn     string
	immh, immb uint32
	q          uint32
}

// newSshr - the Sshr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSshr(b base, q, immh, immb uint32, rd, rn VReg) (Sshr, error) {
	if immh > 0xf || immb > 7 {
		return Sshr{}, fmt.Errorf("arm64.NewSshr: imm out of range")
	}

	return Sshr{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immh: immh,
		immb: immb,
		q:    q,
	}, nil
}

const sshrEnc uint32 = 0x0F000400 // sshr vd, vn, #shift (Q=0 form)

func (i Sshr) ObjDump(_ disasm.ViewCtx) string {
	size, shift := decodeShiftImm(i.immh, i.immb)
	shift = uint32(8<<size) - shift
	return fmt.Sprintf("sshr.%s %s, %s, #0x%x",
		decodeArrangement(i.q, size), i.rd, i.rn, shift)
}

func (i Sshr) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("sshr: %w", err)
	}

	return writeWord(w, sshrEnc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}

// Sshr - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Sshr(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSshr: %w", err)
	}

	if shift == 0 || shift > 8<<size {
		return nil, fmt.Errorf("arm64.NewSshr: shift %d out of range for .%s",
			shift, arr)
	}

	// the stored value is esize-shift, under the size marker bit
	imm := 1<<(3+size) | (8<<size - shift)
	immh, immb := imm>>3, imm&7

	return newSshr(base{}, q, immh, immb, rd, rn)
}

func decodeSshr(w uint32) (Instr, error) {
	in, err := newSshr(
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
