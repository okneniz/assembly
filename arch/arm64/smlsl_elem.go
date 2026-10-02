package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SmlslElem — smlsl{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SmlslElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSmlslElem - the SmlslElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSmlslElem(q, size, idx uint32, rd, rn, rm VReg) (SmlslElem, error) {
	if size == 0 || size > 2 {
		return SmlslElem{}, fmt.Errorf(
			"arm64.NewSmlslElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("SmlslElem", size, idx, rm)
	if err != nil {
		return SmlslElem{}, err
	}

	return SmlslElem{
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
	smlslElemU   uint32 = 0
	smlslElemOpc uint32 = 6
)

func (i SmlslElem) ObjDump(_ disasm.ViewCtx) string {
	name := "smlsl"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i SmlslElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("smlsl: %w", err)
	}

	return writeWord(w, byElemBits(i.q, smlslElemU, i.size, rm, smlslElemOpc, i.idx, rn, rd))
}
