package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Aese — aese.16b vd, vn (the AES round helper; .16b implicit).
type Aese struct {
	base

	rd, rn string
}

// newAese - the Aese constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newAese(b base, rd, rn VReg) (Aese, error) {
	return Aese{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const aeseEnc uint32 = 1311262720 // aese vd, vn

func (i Aese) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("aese.16b %s, %s", i.rd, i.rn)
}

func (i Aese) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("aese: %w", err)
	}

	return writeWord(w, aeseEnc|rd|rn<<5)
}

func (Builder) Aese(rd, rn VReg) (Instr, error) {
	return newAese(base{}, rd, rn)
}

func decodeAese(w uint32) (Instr, error) {
	in, err := newAese(
		newBase(w),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
