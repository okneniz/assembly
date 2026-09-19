package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Uaddw — uaddw{,2}.Arr vd, vn, vm (the widening three-same group:
// the arrangement is the RESULT's - one lane wider than the Rm source;
// Q=1 prints the "2" form).
type Uaddw struct {
	base

	q, size    uint32
	rd, rn, rm string
}

// newUaddw - the Uaddw constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUaddw(b base, q, size uint32, rd, rn, rm VReg) (Uaddw, error) {
	if size > 2 {
		return Uaddw{}, fmt.Errorf(
			"arm64.NewUaddw: the source arrangement is too narrow for .2d results",
		)
	}

	return Uaddw{
		base: b,
		q:    q,
		size: size,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const uaddwEnc uint32 = 773853184 // uaddw vd, vn, vm (Q=0 form)

func (i Uaddw) ObjDump(_ disasm.ViewCtx) string {
	name := "uaddw"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s",
		name, decodeArrangement(i.q, i.size+1), i.rd, i.rn, i.rm)
}

func (i Uaddw) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("uaddw: %w", err)
	}

	return writeWord(w, uaddwEnc|i.q<<30|i.size<<22|rd|rn<<5|rm<<16)
}

func (Builder) Uaddw(rd, rn, rm VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUaddw: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf("arm64.NewUaddw: arrangement too narrow (the source is one lane narrower)")
	}

	return newUaddw(base{}, q, size-1, rd, rn, rm)
}

func decodeUaddw(w uint32) (Instr, error) {
	in, err := newUaddw(
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
