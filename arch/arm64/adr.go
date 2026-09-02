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

// Adr — adr rd, #off: off — the signed byte offset from the address
// of the instruction (the imm21 form, -0x100000..0xfffff). rd — only x
// registers (register 31 reads as zr).
func (Builder) Adr(rd Reg, off int64) (Instr, error) {
	if err := requireClass(
		rd,
		"Adr",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if off < -(1<<20) || off >= 1<<20 {
		return nil, fmt.Errorf(
			"arm64.NewAdr: operand off: %d is out of the imm21 range (-0x100000..0xfffff)",
			off,
		)
	}

	return Adr{rd: rd.name(), off: off}, nil
}

func decodeAdr(w uint32, addr uint64) Instr {
	raw := (w>>5&0x7ffff)<<2 | w>>29&3
	return Adr{
		base: newBase(addr, w),
		rd:   regNameX(w & 0x1f),
		off:  signExtendN(raw, 21),
	}
}

func (i Adr) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("adr %s, #%d", i.rd, i.off)
}

func (i Adr) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, 0x10000000|regBitsX(i.rd)|uint32(i.off&3)<<29|uint32(i.off>>2&0x7ffff)<<5)
}

func (i Adr) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"adr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"PC-relative",
		map[string]any{"Rd": i.rd},
	)
}
