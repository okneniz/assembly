// Package arm64 — per-instruction ARM64 structs: word decoders, formatters
// (objdump notation) and instruction constructors — the Builder methods
// (AddImm, Ldr, ...), the exact inverse of decode.
package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Adc — adc rd, rn, rm (only the 64-bit form is decoded).
type Adc struct {
	base

	rd, rn, rm string
}

// newAdc - the Adc constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newAdc(b base, rd string, rn string, rm string) Adc {
	return Adc{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
	}
}

const adcX uint32 = 0x9A000000

func (i Adc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("adc %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Adc) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("adc: %w", err)
	}

	return writeWord(w, adcX|rd|rn<<5|rm<<16)
}

// Adc — adc rd, rn, rm. Only the 64-bit form; register 31 reads
// as zr (SP/WSP are not allowed — use XZR).
func (Builder) Adc(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(rd, "Adc", "rd",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "Adc", "rn",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	if err := requireClass(rm, "Adc", "rm",
		"register 31 reads as zr — use XZR (only the 64-bit form)", classX, classXZR); err != nil {
		return nil, err
	}

	return newAdc(base{}, rd.name(), rn.name(), rm.name()), nil
}

func decodeAdc(w uint32) Instr {
	return newAdc(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
	)
}
