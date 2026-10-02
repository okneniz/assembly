package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// MlsElem — mls.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type MlsElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newMlsElem - the MlsElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newMlsElem(q, size, idx uint32, rd, rn, rm VReg) (MlsElem, error) {
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
