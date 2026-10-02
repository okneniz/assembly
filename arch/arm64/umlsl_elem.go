package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// UmlslElem — umlsl{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type UmlslElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newUmlslElem - the UmlslElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUmlslElem(b base, q, size, idx uint32, rd, rn, rm VReg) (UmlslElem, error) {
	if size == 0 || size > 2 {
		return UmlslElem{}, fmt.Errorf(
			"arm64.NewUmlslElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("UmlslElem", size, idx, rm)
	if err != nil {
		return UmlslElem{}, err
	}

	return UmlslElem{
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
	umlslElemU   uint32 = 1
	umlslElemOpc uint32 = 6
)

func (i UmlslElem) ObjDump(_ disasm.ViewCtx) string {
	name := "umlsl"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i UmlslElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("umlsl: %w", err)
	}

	return writeWord(w, byElemBits(i.q, umlslElemU, i.size, rm, umlslElemOpc, i.idx, rn, rd))
}
