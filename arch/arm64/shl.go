package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Shl — shl.Arr vd, vn, #shift (the shift lives in immh:immb,
// the plain amount).
type Shl struct {
	rd, rn     string
	immh, immb uint32
	q          uint32
}

// newShl - the Shl constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newShl(q, immh, immb uint32, rd, rn VReg) (Shl, error) {
	if immh > 0xf || immb > 7 {
		return Shl{}, errors.New("arm64.NewShl: imm out of range")
	}

	return Shl{
		rd:   rd.name(),
		rn:   rn.name(),
		immh: immh,
		immb: immb,
		q:    q,
	}, nil
}

const shlEnc uint32 = 0x0F005400 // shl vd, vn, #shift (Q=0 form)

func (i Shl) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("shl: %w", err)
	}

	return writeWord(w, shlEnc|i.q<<30|rd|rn<<5|i.immb<<16|i.immh<<19)
}

func (i Shl) ObjDump(_ disasm.ViewCtx) string {
	size, shift := decodeShiftImm(i.immh, i.immb)
	return fmt.Sprintf("shl.%s %s, %s, #0x%x",
		decodeArrangement(i.q, size), i.rd, i.rn, shift)
}
