package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqdmlslElem — sqdmlsl{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SqdmlslElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSqdmlslElem - the SqdmlslElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqdmlslElem(q, size, idx uint32, rd, rn, rm VReg) (SqdmlslElem, error) {
	if size == 0 || size > 2 {
		return SqdmlslElem{}, fmt.Errorf(
			"arm64.NewSqdmlslElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("SqdmlslElem", size, idx, rm)
	if err != nil {
		return SqdmlslElem{}, err
	}

	return SqdmlslElem{
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
	sqdmlslElemU   uint32 = 0
	sqdmlslElemOpc uint32 = 7
)

func (i SqdmlslElem) ObjDump(_ disasm.ViewCtx) string {
	name := "sqdmlsl"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i SqdmlslElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqdmlsl: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqdmlslElemU, i.size, rm, sqdmlslElemOpc, i.idx, rn, rd))
}
