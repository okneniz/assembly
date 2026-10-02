package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// UmullElem — umull{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type UmullElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newUmullElem - the UmullElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUmullElem(b base, q, size, idx uint32, rd, rn, rm VReg) (UmullElem, error) {
	if size > 2 {
		return UmullElem{}, fmt.Errorf(
			"arm64.NewUmullElem: the source lanes are at most .s (.2d results take .4s sources)",
		)
	}

	err := requireByElemLane("UmullElem", size, idx, rm)
	if err != nil {
		return UmullElem{}, err
	}

	return UmullElem{
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
	umullElemU   uint32 = 1
	umullElemOpc uint32 = 10
)

func (i UmullElem) ObjDump(_ disasm.ViewCtx) string {
	name := "umull"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i UmullElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("umull: %w", err)
	}

	return writeWord(w, byElemBits(i.q, umullElemU, i.size, rm, umullElemOpc, i.idx, rn, rd))
}
