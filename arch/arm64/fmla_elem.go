package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmlaElem — fmla.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmlaElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newFmlaElem - the FmlaElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmlaElem(b base, q, size, idx uint32, rd, rn, rm VReg) (FmlaElem, error) {
	if size != 2 && size != 3 {
		return FmlaElem{}, fmt.Errorf(
			"arm64.NewFmlaElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmlaElem", size, idx, rm)
	if err != nil {
		return FmlaElem{}, err
	}

	return FmlaElem{
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
	fmlaElemU   uint32 = 0
	fmlaElemOpc uint32 = 1
)

func (i FmlaElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmla.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}

func (i FmlaElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmla: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmlaElemU, i.size, rm, fmlaElemOpc, i.idx, rn, rd))
}

// FmlaElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmlaElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmlaElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmlaElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmlaElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmlaElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

func decodeFmlaElem(w uint32) (Instr, error) {
	in, err := newFmlaElem(
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
