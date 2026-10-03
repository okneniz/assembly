package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// UmlalElem — umlal{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type UmlalElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newUmlalElem - the UmlalElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUmlalElem(q, size, idx uint32, rd, rn, rm VReg) (UmlalElem, error) {
	if size == 0 || size > 2 {
		return UmlalElem{}, errors.New(
			"arm64.NewUmlalElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("UmlalElem", size, idx, rm)
	if err != nil {
		return UmlalElem{}, err
	}

	return UmlalElem{
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
	umlalElemU   uint32 = 1
	umlalElemOpc uint32 = 2
)

func (i UmlalElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("umlal: %w", err)
	}

	return writeWord(w, byElemBits(i.q, umlalElemU, i.size, rm, umlalElemOpc, i.idx, rn, rd))
}

func (i UmlalElem) ObjDump(_ disasm.ViewCtx) string {
	name := "umlal"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}
