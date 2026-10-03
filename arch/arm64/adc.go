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
	rd, rn, rm string
}

// newAdc - the Adc constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAdc(rd Reg, rn Reg, rm Reg) (Adc, error) {
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
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const adcX uint32 = 0x9A000000

func (i Adc) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("adc: %w", err)
	}

	return writeWord(w, adcX|rd|rn<<5|rm<<16)
}

func (i Adc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("adc %s, %s, %s", i.rd, i.rn, i.rm)
}
