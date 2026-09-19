package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Saddw — saddw{,2}.Arr vd, vn, vm (the widening three-same group:
// the arrangement is the RESULT's - one lane wider than the Rm source;
// Q=1 prints the "2" form).
type Saddw struct {
	base

	q, size    uint32
	rd, rn, rm string
}

// newSaddw - the Saddw constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSaddw(b base, q, size uint32, rd, rn, rm VReg) (Saddw, error) {
	if size > 2 {
		return Saddw{}, fmt.Errorf(
			"arm64.NewSaddw: the source arrangement is too narrow for .2d results",
		)
	}

	return Saddw{
		base: b,
		q:    q,
		size: size,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const saddwEnc uint32 = 236982272 // saddw vd, vn, vm (Q=0 form)

func (i Saddw) ObjDump(_ disasm.ViewCtx) string {
	name := "saddw"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s",
		name, decodeArrangement(i.q, i.size+1), i.rd, i.rn, i.rm)
}

func (i Saddw) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("saddw: %w", err)
	}

	return writeWord(w, saddwEnc|i.q<<30|i.size<<22|rd|rn<<5|rm<<16)
}

func (Builder) Saddw(rd, rn, rm VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSaddw: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf("arm64.NewSaddw: arrangement too narrow (the source is one lane narrower)")
	}

	return newSaddw(base{}, q, size-1, rd, rn, rm)
}

func decodeSaddw(w uint32) (Instr, error) {
	in, err := newSaddw(
		newBase(w),
		w>>30&1,
		w>>22&3,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(w>>16&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
