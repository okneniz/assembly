package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmlaElem — fmla.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmlaElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newFmlaElem - the FmlaElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmlaElem(b base, q, size, idx uint32, rd, rn, rm VReg) (FmlaElem, error) {
	if size != 2 && size != 3 {
		return FmlaElem{}, fmt.Errorf(
			"arm64.NewFmlaElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmlaElem", size, idx, rm)
	if err != nil {
		return FmlaElem{}, err
	}

	return FmlaElem{
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
	fmlaElemU   uint32 = 0
	fmlaElemOpc uint32 = 1
)

func (i FmlaElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmla.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}

func (i FmlaElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmla: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmlaElemU, i.size, rm, fmlaElemOpc, i.idx, rn, rd))
}
