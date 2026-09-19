package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmulElem — fmul.Arr vd, vn, vm[idx] (the by-element group: a
// vector times one lane of another).
type FmulElem struct {
	base

	q, size, idx uint32
	rd, rn, rm   string
}

// newFmulElem - the FmulElem constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newFmulElem(b base, q, size, idx uint32, rd, rn, rm VReg) (FmulElem, error) {
	if size != 2 && size != 3 {
		return FmulElem{}, fmt.Errorf(
			"arm64.NewFmulElem: only the fp32 (.2s/.4s) and fp64 (.2d) lanes exist",
		)
	}

	err := requireByElemLane("FmulElem", size, idx, rm)
	if err != nil {
		return FmulElem{}, err
	}

	return FmulElem{
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
	fmulElemU   uint32 = 0
	fmulElemOpc uint32 = 9
)

func (i FmulElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, 2) // .2s/.4s
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fmul.%s %s, %s, %s[%d]", arr, i.rd, i.rn, i.rm, i.idx)
}

func (i FmulElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmul: %w", err)
	}

	return writeWord(w, byElemBits(i.q, fmulElemU, i.size, rm, fmulElemOpc, i.idx, rn, rd))
}

// FmulElem - the Builder entry: .2s/.4s (fp32) or .2d (fp64).
func (Builder) FmulElem(rd, rn, rm VReg, arr string, idx uint32) (Instr, error) {
	switch arr {
	case "2s":
		return newFmulElem(base{}, 0, 2, idx, rd, rn, rm)
	case "4s":
		return newFmulElem(base{}, 1, 2, idx, rd, rn, rm)
	case "2d":
		return newFmulElem(base{}, 1, 3, idx, rd, rn, rm)
	}

	return nil, fmt.Errorf("arm64.NewFmulElem: arrangement %q is not one of [2s 4s 2d]", arr)
}

func decodeFmulElem(w uint32) (Instr, error) {
	in, err := newFmulElem(
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
