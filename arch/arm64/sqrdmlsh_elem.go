package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqrdmlshElem — sqrdmlsh.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type SqrdmlshElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSqrdmlshElem - the SqrdmlshElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqrdmlshElem(q, size, idx uint32, rd, rn, rm VReg) (SqrdmlshElem, error) {
	if size == 0 || size == 3 {
		return SqrdmlshElem{}, errors.New(
			"arm64.NewSqrdmlshElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("SqrdmlshElem", size, idx, rm)
	if err != nil {
		return SqrdmlshElem{}, err
	}

	return SqrdmlshElem{
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

func (i SqrdmlshElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqrdmlsh: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqrdmlshElemU, i.size, rm, sqrdmlshElemOpc, i.idx, rn, rd))
}

func (i SqrdmlshElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sqrdmlsh.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}
