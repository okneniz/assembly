package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FcmlaElem — fcmla.Arr vd, vn, vm[idx], #rot (the by-element fused
// multiply-add of complex lanes; the rotation is 0/90/180/270 and
// rides the opcode bits: 0001→#0, 0011→#90, 0101→#180, 0111→#270).
type FcmlaElem struct {
	base

	q, size, idx uint32
	rot          uint32 // 0/90/180/270
	rd, rn, rm   string
}

// newFcmlaElem - the FcmlaElem constructor: validates the operands
// and assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFcmlaElem(b base, q, size, idx, rot uint32, rd, rn, rm VReg) (FcmlaElem, error) {
	err := requireByElemLane("FcmlaElem", size, idx, rm)
	if err != nil {
		return FcmlaElem{}, err
	}

	if rot%90 != 0 || rot > 270 {
		return FcmlaElem{}, fmt.Errorf(
			"arm64.NewFcmlaElem: rotation %d is not one of #0/#90/#180/#270", rot,
		)
	}

	return FcmlaElem{
		base: b,
		q:    q,
		size: size,
		idx:  idx,
		rot:  rot,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

// The opcode bits of the family: U (bit 29) and the opc base (bits
// 15:12) - the rotation adds rot/90*2.
const (
	fcmlaElemU       uint32 = 1
	fcmlaElemOpcBase uint32 = 1
)

func (i FcmlaElem) ObjDump(_ disasm.ViewCtx) string {
	arr := decodeArrangement(i.q, i.size)
	if i.size == 3 {
		arr = "2d"
	}

	return fmt.Sprintf("fcmla.%s %s, %s, %s[%d], #%d",
		arr, i.rd, i.rn, i.rm, i.idx, i.rot)
}

func (i FcmlaElem) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fcmla: %w", err)
	}

	opc := fcmlaElemOpcBase + i.rot/90*2
	return writeWord(w, byElemBits(i.q, fcmlaElemU, i.size, rm, opc, i.idx, rn, rd))
}

// FcmlaElem - the Builder entry: .4h/.8h/.2s/.4s/.2d lanes.
func (Builder) FcmlaElem(rd, rn, rm VReg, arr string, idx, rot uint32) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewFcmlaElem: %w", err)
	}

	if size == 0 {
		return nil, fmt.Errorf(
			"arm64.NewFcmlaElem: arrangement %q is not one of the complex lane widths", arr,
		)
	}

	return newFcmlaElem(base{}, q, size, idx, rot, rd, rn, rm)
}

func decodeFcmlaElem(w uint32) (Instr, error) {
	size := w >> 22 & 3
	in, err := newFcmlaElem(
		newBase(w),
		w>>30&1,
		size,
		byElemIndex(w, size),
		w>>12&0xf/2*90%360,
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
		newVReg(uint8(byElemVm(w, size))),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
