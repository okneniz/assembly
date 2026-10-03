package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmulxElem — fmulx.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmulxElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newFmulxElem - the FmulxElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmulxElem(q, size, idx uint32, rd, rn, rm VReg) (FmulxElem, error) {
	if size != 2 && size != 3 {
		return FmulxElem{}, errors.New(
			"arm64.NewFmulxElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmulxElem", size, idx, rm)
	if err != nil {
		return FmulxElem{}, err
	}

	return FmulxElem{
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
	fmulxElemU   uint32 = 1
	fmulxElemOpc uint32 = 9
)

func (i FmulxElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmulx: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmulxElemU, i.size, rm, fmulxElemOpc, i.idx, rn, rd))
}

func (i FmulxElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmulx.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}
