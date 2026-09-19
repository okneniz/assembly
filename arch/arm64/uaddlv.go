package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Uaddlv — uaddlv.Arr hN/sN/dN, vn (scalar dest by size).
type Uaddlv struct {
	base

	rd, rn  string // rd - the scalar print name (hN/sN/dN)
	q, size uint32
}

// newUaddlv - the Uaddlv constructor: the struct is assembled only
// here (the Builder method and the decoder call it); the destination
// print name is the scalar view of the register number by size.
func newUaddlv(b base, q, size uint32, rd, rn VReg) (Uaddlv, error) {
	scalar := fmt.Sprintf("d%d", rd.Num())
	if size == 0 {
		scalar = fmt.Sprintf("h%d", rd.Num())
	} else if size == 1 {
		scalar = fmt.Sprintf("s%d", rd.Num())
	}

	return Uaddlv{
		base: b,
		rd:   scalar,
		rn:   rn.name(),
		q:    q,
		size: size,
	}, nil
}

// Uaddlv - the Builder entry: .8b/.16b/.4h/.8h lanes.
func (Builder) Uaddlv(rd, rn VReg, arr string) (Instr, error) {
	q, size, err := arrBits(arr)
	if err != nil {
		return nil, fmt.Errorf("arm64.NewUaddlv: %w", err)
	}

	if size > 1 {
		return nil, fmt.Errorf("arm64.NewUaddlv: arrangement %q is not one of [8b 16b 4h 8h]", arr)
	}

	return newUaddlv(base{}, q, size, rd, rn)
}

const uaddlvEnc uint32 = 0x2E303800

func (i Uaddlv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("uaddlv.%s %s, %s", decodeArrangement(i.q, i.size), i.rd, i.rn)
}

func (i Uaddlv) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(fmt.Sprintf("v%d", regIndex(i.rd)), i.rn)
	if err != nil {
		return 0, fmt.Errorf("uaddlv: %w", err)
	}

	return writeWord(w, uaddlvEnc|rd|rn<<5)
}

func decodeUaddlv(w uint32) (Instr, error) {
	rd := vReg(w & 0x1f)
	size := w >> 22 & 3
	switch size {
	case 0:
		rd = fmt.Sprintf("h%d", regIndex(rd))
	case 1:
		rd = fmt.Sprintf("s%d", regIndex(rd))
	}

	return Uaddlv{
		base: newBase(w),
		rd:   rd,
		rn:   vReg(w >> 5 & 0x1f),
		q:    w >> 30 & 1,
		size: size,
	}, nil
}
