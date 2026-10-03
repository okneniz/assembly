package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MlaElem — mla.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type MlaElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newMlaElem - the MlaElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newMlaElem(q, size, idx uint32, rd, rn, rm VReg) (MlaElem, error) {
	if size == 0 || size == 3 {
		return MlaElem{}, errors.New("arm64.NewMlaElem: only the .h and .s integer lanes exist")
	}

	err := requireByElemLane("MlaElem", size, idx, rm)
	if err != nil {
		return MlaElem{}, err
	}

	return MlaElem{
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
	mlaElemU   uint32 = 1
	mlaElemOpc uint32 = 0
)

func (i MlaElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("mla: %w", err)
	}

	return writeWord(w, byElemBits(i.q, mlaElemU, i.size, rm, mlaElemOpc, i.idx, rn, rd))
}

func (i MlaElem) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mla.%s %s, %s, %s[%d]",
		decodeArrangement(i.q, i.size), i.rd, i.rn, i.rm, i.idx)
}
