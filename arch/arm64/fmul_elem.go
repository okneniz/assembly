package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmulElem — fmul.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmulElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newFmulElem - the FmulElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmulElem(b base, q, size, idx uint32, rd, rn, rm VReg) (FmulElem, error) {
	if size != 2 && size != 3 {
		return FmulElem{}, fmt.Errorf(
			"arm64.NewFmulElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmulElem", size, idx, rm)
	if err != nil {
		return FmulElem{}, err
	}

	return FmulElem{
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
	fmulElemU   uint32 = 0
	fmulElemOpc uint32 = 9
)

func (i FmulElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmul.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}

func (i FmulElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmul: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmulElemU, i.size, rm, fmulElemOpc, i.idx, rn, rd))
}
