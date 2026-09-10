package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Adr — adr rd, #off (imm21: immhi:immlo, byte offset).
type Adr struct {
	base

	rd  string
	off int64
}

// newAdr - the Adr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAdr(b base, rd Reg, off int64) (Adr, error) {
	err := requireClass(
		rd,
		"Adr",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	)

	if err != nil {
		return Adr{}, err
	}

	if off < -(1<<20) || off >= 1<<20 {
		return Adr{}, fmt.Errorf(
			"arm64.NewAdr: operand off: %d is out of the imm21 range (-0x100000..0xfffff)",
			off,
		)
	}

	return Adr{
		base: b,
		rd:   rd.name(),
		off:  off,
	}, nil
}

func (i Adr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("adr %s, #%d", i.rd, i.off)
}

func (i Adr) Encode(w io.Writer) (int64, error) {
	return writeWord(w, 0x10000000|regBitsX(i.rd)|uint32(i.off&3)<<29|uint32(i.off>>2&0x7ffff)<<5)
}

// Adr — adr rd, #off: off — the signed byte offset from the address
// of the instruction (the imm21 form, -0x100000..0xfffff). rd — only x
// registers (register 31 reads as zr).
func (Builder) Adr(rd Reg, off int64) (Instr, error) {
	return newAdr(base{}, rd, off)
}

func decodeAdr(w uint32) (Instr, error) {
	raw := (w>>5&0x7ffff)<<2 | w>>29&3
	in, err := newAdr(
		newBase(w),
		gprOf(w&0x1f, true),
		signExtendN(raw, 21),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
