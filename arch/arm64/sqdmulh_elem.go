package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqdmulhElem — sqdmulh.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type SqdmulhElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newSqdmulhElem - the SqdmulhElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqdmulhElem(b base, q, size, idx uint32, rd, rn, rm VReg) (SqdmulhElem, error) {
	if size == 0 || size == 3 {
		return SqdmulhElem{}, fmt.Errorf(
			"arm64.NewSqdmulhElem: only the .h and .s integer lanes exist",
		)
	}

	err := requireByElemLane("SqdmulhElem", size, idx, rm)
	if err != nil {
		return SqdmulhElem{}, err
	}

	return SqdmulhElem{
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
	sqdmulhElemU   uint32 = 0
	sqdmulhElemOpc uint32 = 12
)

func (i SqdmulhElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sqdmulh.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}

func (i SqdmulhElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqdmulh: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqdmulhElemU, i.size, rm, sqdmulhElemOpc, i.idx, rn, rd))
}
