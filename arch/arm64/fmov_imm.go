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

func decodeFmovImmOf(isS bool, enc uint32, rdK fpKind) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		imm8 := w >> 13 & 0xff
		rd := fpReg(w&0x1f, rdK)
		if isS {
			v := vfpExpandImm32(imm8)
			return FmovImm{
				base: newBase(addr, w),
				rd:   rd,
				val:  float64(v),
				text: fmt.Sprintf("%.8f", v),
				isS:  true,
				enc:  enc,
				rdK:  rdK,
			}
		}

		v := vfpExpandImm64(imm8)
		return FmovImm{
			base: newBase(addr, w),
			rd:   rd,
			val:  v,
			text: fmt.Sprintf("%.8f", v),
			enc:  enc,
			rdK:  rdK,
		}
	}
}

func (i FmovImm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmov %s, #%s", i.rd, i.text)
}

func (i FmovImm) Encode(w io.Writer, pc uint64) (int64, error) {
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

func (i FmovImm) MarshalJSON() ([]byte, error) {
	return i.marshal("fmov", i.ObjDump(disasm.DefaultViewCtx()), "FP", map[string]any{"Rd": i.rd})
}
