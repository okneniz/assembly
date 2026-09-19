package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SqdmullElem — sqdmull{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SqdmullElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newSqdmullElem - the SqdmullElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqdmullElem(b base, q, size, idx uint32, rd, rn, rm VReg) (SqdmullElem, error) {
	if size > 2 {
		return SqdmullElem{}, fmt.Errorf(
			"arm64.NewSqdmullElem: the source lanes are at most .s (.2d results take .4s sources)",
		)
	}

	err := requireByElemLane("SqdmullElem", size, idx, rm)
	if err != nil {
		return SqdmullElem{}, err
	}

	return SqdmullElem{
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

// SqdmullElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SqdmullElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "8h", "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSqdmullElem: arrangement %q is not one of [8h 4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSqdmullElem: %w", err)
	}

	return newSqdmullElem(base{}, q, size-1, idx, rd, rn, rm)
}

func decodeSqdmullElem(w uint32) (Instr, error) {
	in, err := newSqdmullElem(
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
