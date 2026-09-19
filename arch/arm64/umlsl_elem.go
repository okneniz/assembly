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
	if size > 2 {
		return UmlslElem{}, fmt.Errorf(
			"arm64.NewUmlslElem: the source lanes are at most .s (.2d results take .4s sources)",
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

// UmlslElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) UmlslElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "8h", "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewUmlslElem: arrangement %q is not one of [8h 4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUmlslElem: %w", err)
	}

	return newUmlslElem(base{}, q, size-1, idx, rd, rn, rm)
}

func decodeUmlslElem(w uint32) (Instr, error) {
	in, err := newUmlslElem(
		newBase(w),
		w>>30&1,
		w>>22&3,
		byElemIndex(w, w>>22&3),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, w>>22&3))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
