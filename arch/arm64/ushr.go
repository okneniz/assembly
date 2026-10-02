package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ushr — ushr.Arr vd, vn, #shift (the shift lives in immh:immb,
// the stored value is esize-shift (printed shifted right)).
type Ushr struct {
	base

	rd, rn     string
	immh, immb uint32
	q          uint32
}

// newUshr - the Ushr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUshr(b base, q, immh, immb uint32, rd, rn VReg) (Ushr, error) {
	if immh > 0xf || immb > 7 {
		return Ushr{}, fmt.Errorf("arm64.NewUshr: imm out of range")
	}

	return Ushr{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		immh: immh,
		immb: immb,
		q:    q,
	}, nil
}

const ushrEnc uint32 = 0x2F000400 // ushr vd, vn, #shift (Q=0 form)

func (i Ushr) ObjDump(_ disasm.ViewCtx) string {
	size, shift := decodeShiftImm(i.immh, i.immb)
	shift = uint32(8<<size) - shift
	return fmt.Sprintf("ushr.%s %s, %s, #0x%x",
		decodeArrangement(i.q, size), i.rd, i.rn, shift)
}

func (i Ushr) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("ushr: %w", err)
	}

	return writeWord(w, ushrEnc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}
