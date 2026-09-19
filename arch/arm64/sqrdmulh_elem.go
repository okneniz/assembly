package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqrdmulhElem — sqrdmulh.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type SqrdmulhElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newSqrdmulhElem - the SqrdmulhElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqrdmulhElem(b base, q, size, idx uint32, rd, rn, rm VReg) (SqrdmulhElem, error) {
	if size == 0 || size == 3 {
		return SqrdmulhElem{}, fmt.Errorf(
			"arm64.NewSqrdmulhElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("SqrdmulhElem", size, idx, rm)
	if err != nil {
		return SqrdmulhElem{}, err
	}

	return SqrdmulhElem{
		base: b,
		q:    q,
		size: size,
		idx:  idx,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

// The opcode bits of the family: U (bit 29) and opc (bits 15:12).
const (
	sqrdmulhElemU   uint32 = 0
	sqrdmulhElemOpc uint32 = 13
)

func (i SqrdmulhElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sqrdmulh.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}

func (i SqrdmulhElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqrdmulh: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqrdmulhElemU, i.size, rm, sqrdmulhElemOpc, i.idx, rn, rd))
}

// SqrdmulhElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) SqrdmulhElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqrdmulhElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf("arm64.NewSqrdmulhElem: arrangement %q is not one of the integer lane widths", arr)
	}

	return newSqrdmulhElem(base{}, q, size, idx, rd, rn, rm)
}

func decodeSqrdmulhElem(w uint32) (Instr, error) {
	in, err := newSqrdmulhElem(
		newBase(w),
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
