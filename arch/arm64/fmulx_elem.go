package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmulxElem — fmulx.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmulxElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newFmulxElem - the FmulxElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmulxElem(b base, q, size, idx uint32, rd, rn, rm VReg) (FmulxElem, error) {
	if size != 2 && size != 3 {
		return FmulxElem{}, fmt.Errorf(
			"arm64.NewFmulxElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmulxElem", size, idx, rm)
	if err != nil {
		return FmulxElem{}, err
	}

	return FmulxElem{
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
	fmulxElemU   uint32 = 1
	fmulxElemOpc uint32 = 9
)

func (i FmulxElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmulx.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}

func (i FmulxElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmulx: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmulxElemU, i.size, rm, fmulxElemOpc, i.idx, rn, rd))
}

// FmulxElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmulxElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmulxElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmulxElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmulxElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmulxElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

func decodeFmulxElem(w uint32) (Instr, error) {
	in, err := newFmulxElem(
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
