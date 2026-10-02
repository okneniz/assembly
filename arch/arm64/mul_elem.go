package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MulElem — mul.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type MulElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newMulElem - the MulElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newMulElem(b base, q, size, idx uint32, rd, rn, rm VReg) (MulElem, error) {
	if size == 0 || size == 3 {
		return MulElem{}, fmt.Errorf(
			"arm64.NewMulElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("MulElem", size, idx, rm)
	if err != nil {
		return MulElem{}, err
	}

	return MulElem{
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
	mulElemU   uint32 = 0
	mulElemOpc uint32 = 8
)

func (i MulElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mul.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}

func (i MulElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("mul: %w", err)
	}

	return writeWord(w, byElemBits(i.q, mulElemU, i.size, rm, mulElemOpc, i.idx, rn, rd))
}
