package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Shl — shl.Arr vd, vn, #shift (the shift lives in immh:immb,
// the plain amount).
type Shl struct {
	base

	rd, rn     string
	immh, immb uint32
	q          uint32
}

// newShl - the Shl constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newShl(b base, q, immh, immb uint32, rd, rn VReg) (Shl, error) {
	if immh > 0xf || immb > 7 {
		return Shl{}, fmt.Errorf("arm64.NewShl: imm out of range")
	}

	return Shl{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immh: immh,
		immb: immb,
		q:    q,
	}, nil
}

const shlEnc uint32 = 0x0F005400 // shl vd, vn, #shift (Q=0 form)

func (i Shl) ObjDump(_ disasm.ViewCtx) string {
	size, shift := decodeShiftImm(i.immh, i.immb)
	return fmt.Sprintf("shl.%s %s, %s, #0x%x",
		decodeArrangement(i.q, size), i.rd, i.rn, shift)
}

func (i Shl) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("shl: %w", err)
	}

	return writeWord(w, shlEnc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}

// Shl - the Builder entry: the shift as the written amount on the
// arrangement's lanes (the raw immh:immb stays a decode detail).
func (Builder) Shl(rd, rn VReg, arr string, shift uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewShl: %w", err)
	}

	if shift >= 8<<size {
		return nil, fmt.Errorf("arm64.NewShl: shift %d out of range for .%s",
			shift, arr)
	}

	// the immediate field carries the size marker above the shift
	imm := 1<<(3+size) | shift
	immh, immb := imm>>3, imm&7

	return newShl(base{}, q, immh, immb, rd, rn)
}

func decodeShl(w uint32) (Instr, error) {
	in, err := newShl(
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
