package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqdmullElem — sqdmull{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SqdmullElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSqdmullElem - the SqdmullElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqdmullElem(q, size, idx uint32, rd, rn, rm VReg) (SqdmullElem, error) {
	if size == 0 || size > 2 {
		return SqdmullElem{}, fmt.Errorf(
			"arm64.NewSqdmullElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("SqdmullElem", size, idx, rm)
	if err != nil {
		return SqdmullElem{}, err
	}

	return SqdmullElem{
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
	sqdmullElemU   uint32 = 0
	sqdmullElemOpc uint32 = 11
)

func (i SqdmullElem) ObjDump(_ disasm.ViewCtx) string {
	name := "sqdmull"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i SqdmullElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqdmull: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqdmullElemU, i.size, rm, sqdmullElemOpc, i.idx, rn, rd))
}
