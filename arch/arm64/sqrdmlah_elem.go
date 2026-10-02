package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqrdmlahElem — sqrdmlah.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type SqrdmlahElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSqrdmlahElem - the SqrdmlahElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqrdmlahElem(q, size, idx uint32, rd, rn, rm VReg) (SqrdmlahElem, error) {
	if size == 0 || size == 3 {
		return SqrdmlahElem{}, fmt.Errorf(
			"arm64.NewSqrdmlahElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("SqrdmlahElem", size, idx, rm)
	if err != nil {
		return SqrdmlahElem{}, err
	}

	return SqrdmlahElem{
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
	sqrdmlahElemU   uint32 = 1
	sqrdmlahElemOpc uint32 = 13
)

func (i SqrdmlahElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sqrdmlah.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}

func (i SqrdmlahElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqrdmlah: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqrdmlahElemU, i.size, rm, sqrdmlahElemOpc, i.idx, rn, rd))
}
