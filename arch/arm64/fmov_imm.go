package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// FmovImm — fmov d/s rd, #imm8 (vfpExpandImm64/32, #%.8f format).
type FmovImm struct {
	base

	rd   string
	val  float64
	text string
	isS  bool
	enc  uint32
	rdK  fpKind
}

// newFmovImm - the FmovImm constructor: assembles the struct (the
// Builder method delegates here; the decoder calls it with values read
// from the word). The imm8 search stays in Encode: it is shared with
// the decode path, which stores the expanded value's text.
func newFmovImm(b base, rd FReg, val float64, text string) (FmovImm, error) {
	return FmovImm{
		base: b,
		rd:   rd.name(),
		val:  val,
		text: text,
		isS:  !rd.Is64(),
		enc:  fmovImmEnc(rd),
		rdK:  rd.kind(),
	}, nil
}

const (
	fmovImmDEnc uint32 = 0x1E601000 // fmov dd, #imm
	fmovImmSEnc uint32 = 0x1E201000 // fmov sd, #imm
)

// fmovImmEnc — the immediate-form encoding by the destination kind.
func fmovImmEnc(rd FReg) uint32 {
	if rd.Is64() {
		return fmovImmDEnc
	}

	return fmovImmSEnc
}

func (Builder) FmovImm(rd FReg, val float64) (Instr, error) {
	return newFmovImm(base{}, rd, val, fmt.Sprintf("%.8f", val))
}

func decodeFmovImmOf(isS bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		imm8 := w >> 13 & 0xff
		rd := newFReg(uint8(w&0x1f), !isS)
		if isS {
			v := vfpExpandImm32(imm8)
			return newFmovImm(newBase(w), rd, float64(v), fmt.Sprintf("%.8f", v))
		}

		v := vfpExpandImm64(imm8)
		return newFmovImm(newBase(w), rd, v, fmt.Sprintf("%.8f", v))
	}
}

func (i FmovImm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmov %s, #%s", i.rd, i.text)
}

func (i FmovImm) Encode(w io.Writer) (int64, error) {
	rd, err := armRegNum(i.rd)
	if err != nil {
		return 0, fmt.Errorf("fmov: %w", err)
	}

	// text → imm8: enumeration of 256 values (canonical reverse search)
	for imm8 := range uint32(256) {
		if i.isS {
			if fmt.Sprintf("%.8f", vfpExpandImm32(imm8)) == i.text {
				return writeWord(w, i.enc|rd|imm8<<13)
			}
		} else if fmt.Sprintf("%.8f", vfpExpandImm64(imm8)) == i.text {
			return writeWord(w, i.enc|rd|imm8<<13)
		}
	}

	return 0, errors.New("fmov: bad imm")
}
