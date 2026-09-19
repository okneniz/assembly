package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sri — sri.Arr vd, vn, #shift (the shift lives in immh:immb,
// the stored value is esize-shift (printed shifted right)).
type Sri struct {
	base

	rd, rn     string
	immh, immb uint32
	q          uint32
}

// newSri - the Sri constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSri(b base, q, immh, immb uint32, rd, rn VReg) (Sri, error) {
	if immh > 0xf || immb > 7 {
		return Sri{}, fmt.Errorf("arm64.NewSri: imm out of range")
	}

	return Sri{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immh: immh,
		immb: immb,
		q:    q,
	}, nil
}

const sriEnc uint32 = 0x2F004400 // sri vd, vn, #shift (Q=0 form)

func (i Sri) ObjDump(_ disasm.ViewCtx) string {
	size, shift := decodeShiftImm(i.immh, i.immb)
	shift = uint32(8<<size) - shift
	return fmt.Sprintf("sri.%s %s, %s, #0x%x",
		decodeArrangement(i.q, size), i.rd, i.rn, shift)
}

func (i Sri) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("sri: %w", err)
	}

	return writeWord(w, sriEnc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}

// Sri - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Sri(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSri: %w", err)
	}

	if shift == 0 || shift > 8<<size {
		return nil, fmt.Errorf("arm64.NewSri: shift %d out of range for .%s",
			shift, arr)
	}

	// the stored value is esize-shift, under the size marker bit
	imm := 1<<(3+size) | (8<<size - shift)
	immh, immb := imm>>3, imm&7

	return newSri(base{}, q, immh, immb, rd, rn)
}

func decodeSri(w uint32) (Instr, error) {
	in, err := newSri(
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
