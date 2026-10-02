package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SmullElem — smull{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SmullElem struct {
	q, size, idx uint32
	rd, rn, rm   string
}

// newSmullElem - the SmullElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSmullElem(q, size, idx uint32, rd, rn, rm VReg) (SmullElem, error) {
	if size == 0 || size > 2 {
		return SmullElem{}, fmt.Errorf(
			"arm64.NewSmullElem: only the .h and .s source lanes exist (the .8h-result class is unallocated)",
		)
	}

	err := requireByElemLane("SmullElem", size, idx, rm)
	if err != nil {
		return SmullElem{}, err
	}

	return SmullElem{
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
	smullElemU   uint32 = 0
	smullElemOpc uint32 = 10
)

func (i SmullElem) ObjDump(_ disasm.ViewCtx) string {
	name := "smull"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i SmullElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("smull: %w", err)
	}

	return writeWord(w, byElemBits(i.q, smullElemU, i.size, rm, smullElemOpc, i.idx, rn, rd))
}
