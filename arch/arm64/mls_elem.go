package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MlsElem — mls.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type MlsElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newMlsElem - the MlsElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newMlsElem(b base, q, size, idx uint32, rd, rn, rm VReg) (MlsElem, error) {
	if size == 0 || size == 3 {
		return MlsElem{}, fmt.Errorf(
			"arm64.NewMlsElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("MlsElem", size, idx, rm)
	if err != nil {
		return MlsElem{}, err
	}

	return MlsElem{
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
	mlsElemU   uint32 = 1
	mlsElemOpc uint32 = 4
)

func (i MlsElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mls.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}

func (i MlsElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("mls: %w", err)
	}

	return writeWord(w, byElemBits(i.q, mlsElemU, i.size, rm, mlsElemOpc, i.idx, rn, rd))
}

// MlsElem - the Builder entry: the arrangement is the printed one
// (.4h/.8h/.2s/.4s).
func (Builder) MlsElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewMlsElem: %w", err)
	}

	if size == 0 || size == 3 {
		return nil, fmt.Errorf("arm64.NewMlsElem: arrangement %q is not one of the integer lane widths", arr)
	}

	return newMlsElem(base{}, q, size, idx, rd, rn, rm)
}

func decodeMlsElem(w uint32) (Instr, error) {
	in, err := newMlsElem(
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
