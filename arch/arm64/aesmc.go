package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Aesmc — aesmc.16b vd, vn (the AES round helper; .16b implicit).
type Aesmc struct {
	base

	rd, rn string
}

// newAesmc - the Aesmc constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newAesmc(b base, rd, rn VReg) (Aesmc, error) {
	return Aesmc{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
	}, nil
}

const aesmcEnc uint32 = 1311270912 // aesmc vd, vn

func (i Aesmc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("aesmc.16b %s, %s", i.rd, i.rn)
}

func (i Aesmc) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("aesmc: %w", err)
	}

	return writeWord(w, aesmcEnc|rd|rn<<5)
}

func (Builder) Aesmc(rd, rn VReg) (Instr, error) {
	return newAesmc(base{}, rd, rn)
}

func decodeAesmc(w uint32) (Instr, error) {
	in, err := newAesmc(
		newBase(w),
		newVReg(uint8(w&0x1f)),
		newVReg(uint8(w>>5&0x1f)),
	)
	if err != nil {
		return decodeUnknown(w) // unencodable operand bits: data
	}

	return in, nil
}
