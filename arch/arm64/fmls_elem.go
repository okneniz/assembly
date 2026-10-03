package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmlsElem — fmls.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmlsElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newFmlsElem - the FmlsElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmlsElem(q, size, idx uint32, rd, rn, rm VReg) (FmlsElem, error) {
	if size != 2 && size != 3 {
		return FmlsElem{}, errors.New(
			"arm64.NewFmlsElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmlsElem", size, idx, rm)
	if err != nil {
		return FmlsElem{}, err
	}

	return FmlsElem{
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
	fmlsElemU   uint32 = 0
	fmlsElemOpc uint32 = 5
)

func (i FmlsElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmls: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmlsElemU, i.size, rm, fmlsElemOpc, i.idx, rn, rd))
}

func (i FmlsElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmls.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}
