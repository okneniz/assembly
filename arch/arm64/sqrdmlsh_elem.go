package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqrdmlshElem — sqrdmlsh.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type SqrdmlshElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newSqrdmlshElem - the SqrdmlshElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqrdmlshElem(b base, q, size, idx uint32, rd, rn, rm VReg) (SqrdmlshElem, error) {
	if size == 0 || size == 3 {
		return SqrdmlshElem{}, fmt.Errorf(
			"arm64.NewSqrdmlshElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("SqrdmlshElem", size, idx, rm)
	if err != nil {
		return SqrdmlshElem{}, err
	}

	return SqrdmlshElem{
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
	sqrdmlshElemU   uint32 = 1
	sqrdmlshElemOpc uint32 = 15
)

func (i SqrdmlshElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sqrdmlsh.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}

func (i SqrdmlshElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqrdmlsh: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqrdmlshElemU, i.size, rm, sqrdmlshElemOpc, i.idx, rn, rd))
}

// SqrdmlshElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) SqrdmlshElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqrdmlshElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf("arm64.NewSqrdmlshElem: arrangement %q is not one of the integer lane widths", arr)
	}

	return newSqrdmlshElem(base{}, q, size, idx, rd, rn, rm)
}

func decodeSqrdmlshElem(w uint32) (Instr, error) {
	in, err := newSqrdmlshElem(
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
