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
	q, size, idx uint32
	rot          uint32 // 0/90/180/270
	rd, rn, rm   string
}

// newFcmlaElem - the FcmlaElem constructor: validates the operands
// and assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word). Unlike the rest of
// the by-element group the Vm field is always 5 bits wide and the lane
// index lives in {L (b21), H (b11)}: .4h/.4s carry one bit, .8h two.
func newFcmlaElem(q, size, idx, rot uint32, rd, rn, rm VReg) (FcmlaElem, error) {
	var maxIdx uint32
	switch {
	case size == 1 && q == 0: // .4h
		maxIdx = 1
	case size == 1 && q == 1: // .8h
		maxIdx = 3
	case size == 2 && q == 1: // .4s
		maxIdx = 1
	default: // .2s/.2d/.b/.d — unallocated (clang refuses, word decodes unknown)
		return FcmlaElem{}, fmt.Errorf(
			"arm64.NewFcmlaElem: arrangement %q is not one of [4h 8h 4s]",
			decodeArrangement(q, size),
		)
	}

	if idx > maxIdx {
		return FcmlaElem{}, fmt.Errorf(
			"arm64.NewFcmlaElem: lane index %d out of range (0..%d)", idx, maxIdx,
		)
	}

	if rot%90 != 0 || rot > 270 {
		return FcmlaElem{}, fmt.Errorf(
			"arm64.NewFcmlaElem: rotation %d is not one of #0/#90/#180/#270", rot,
		)
	}

	return FcmlaElem{
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

	// the lane index: .4h — L (b21) is the low bit, H (b11) the high
	// one; .4s — H only (clang: the .4s lane range is [0, 1] and the
	// set bit is b11)
	var idxBits uint32
	if i.size == 1 {
		idxBits = i.idx&1<<21 | i.idx>>1<<11
	} else {
		idxBits = i.idx << 11
	}

	opc := fcmlaElemOpcBase + i.rot/90*2
	return writeWord(w, i.q<<30|fcmlaElemU<<29|0x0F<<24|i.size<<22|
		rm<<16|opc<<12|idxBits|rn<<5|rd)
}
