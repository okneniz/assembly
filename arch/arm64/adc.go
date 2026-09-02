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

const adcX uint32 = 0x9A000000

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

	return Adc{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

func decodeAdc(w uint32, addr uint64) Instr {
	return Adc{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i Adc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("adc %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Adc) Encode(w io.Writer, pc uint64) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("adc: %w", err)
	}

	return writeWord(w, adcX|rd|rn<<5|rm<<16)
}

func (i Adc) MarshalJSON() ([]byte, error) {
	return i.marshal("adc", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm})
}
