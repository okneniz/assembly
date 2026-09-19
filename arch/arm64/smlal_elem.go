package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SmlalElem — smlal{,2}.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another; the result lanes are one width wider).
type SmlalElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newSmlalElem - the SmlalElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSmlalElem(b base, q, size, idx uint32, rd, rn, rm VReg) (SmlalElem, error) {
	if size > 2 {
		return SmlalElem{}, fmt.Errorf(
			"arm64.NewSmlalElem: the source lanes are at most .s (.2d results take .4s sources)",
		)
	}

	err := requireByElemLane("SmlalElem", size, idx, rm)
	if err != nil {
		return SmlalElem{}, err
	}

	return SmlalElem{
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
	smlalElemU   uint32 = 0
	smlalElemOpc uint32 = 2
)

func (i SmlalElem) ObjDump(_ disasm.ViewCtx) string {
	name := "smlal"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s[%d]",
		name, decodeArrangement(1, i.size+1), i.rd, i.rn, i.rm, i.idx)
}

func (i SmlalElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("smlal: %w", err)
	}

	return writeWord(w, byElemBits(i.q, smlalElemU, i.size, rm, smlalElemOpc, i.idx, rn, rd))
}

// SmlalElem - the Builder entry: the arrangement is the RESULT's
// (.8h/.4s/.2d); two selects the "2" form (the upper halves).
func (Builder) SmlalElem(rd, rn, rm VReg, arr string, two bool, idx uint32) (Instr, error) {
	q := uint32(0)
	if two {
		q = 1
	}

	switch arr {
	case "8h", "4s", "2d":
	default:
		return nil, fmt.Errorf("arm64.NewSmlalElem: arrangement %q is not one of [8h 4s 2d]", arr)
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewSmlalElem: %w", err)
	}

	return newSmlalElem(base{}, q, size-1, idx, rd, rn, rm)
}

func decodeSmlalElem(w uint32) (Instr, error) {
	in, err := newSmlalElem(
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
