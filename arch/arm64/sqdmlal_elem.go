package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqdmlalElem — sqdmlal{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SqdmlalElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSqdmlalElem - the SqdmlalElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqdmlalElem(q, size, idx uint32, rd, rn, rm VReg) (SqdmlalElem, error) {
	if size == 0 || size > 2 {
		return SqdmlalElem{}, fmt.Errorf(
			"arm64.NewSqdmlalElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("SqdmlalElem", size, idx, rm)
	if err != nil {
		return SqdmlalElem{}, err
	}

	return SqdmlalElem{
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
	sqdmlalElemU   uint32 = 0
	sqdmlalElemOpc uint32 = 3
)

func (i SqdmlalElem) ObjDump(_ disasm.ViewCtx) string {
	name := "sqdmlal"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i SqdmlalElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqdmlal: %w", err)
	}

	return writeWord(w, byElemBits(i.q, sqdmlalElemU, i.size, rm, sqdmlalElemOpc, i.idx, rn, rd))
}
