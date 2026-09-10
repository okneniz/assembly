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

// newAdc - the Adc constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAdc(b base, rd Reg, rn Reg, rm Reg) (Adc, error) {
	err := requireClass(
		rd,
		"Adc",
		"rd",
		"register 31 reads as zr — use XZR (only the 64-bit form)",
		classX,
		classXZR,
	)

	if err != nil {
		return Adc{}, err
	}

	err = requireClass(
		rn,
		"Adc",
		"rn",
		"register 31 reads as zr — use XZR (only the 64-bit form)",
		classX,
		classXZR,
	)

	if err != nil {
		return Adc{}, err
	}

	err = requireClass(
		rm,
		"Adc",
		"rm",
		"register 31 reads as zr — use XZR (only the 64-bit form)",
		classX,
		classXZR,
	)

	if err != nil {
		return Adc{}, err
	}

	return Adc{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
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

func (Builder) Adc(rd, rn, rm Reg) (Instr, error) {
	return newAdc(base{}, rd, rn, rm)
}

func decodeAdc(w uint32) (Instr, error) {
	in, err := newAdc(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		gprOf(w>>16&0x1f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
