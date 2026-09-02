package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Adrp — adrp rd, #imm21 ; #page (absolute page annotation).
type Adrp struct {
	base

	rd   string
	off  int64
	page uint64
}

// Adrp — adrp rd, #off: off — the signed imm21 count of 4KB pages
// from the page of the instruction (-0x100000..0xfffff pages). The
// absolute-page annotation of the decoded form needs the instruction
// address and stays zero here. rd — only x registers (register 31 reads
// as zr).
func (Builder) Adrp(rd Reg, off int64) (Instr, error) {
	if err := requireClass(
		rd,
		"Adrp",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if off < -(1<<20) || off >= 1<<20 {
		return nil, fmt.Errorf(
			"arm64.NewAdrp: operand off: %d is out of the imm21 page range (-0x100000..0xfffff pages)",
			off,
		)
	}

	return Adrp{rd: rd.name(), off: off}, nil
}

func decodeAdrp(w uint32, addr uint64) Instr {
	raw := (w>>5&0x7ffff)<<2 | w>>29&3
	imm21 := signExtendN(raw, 21)
	page := (int64(addr) & ^int64(0xFFF)) + imm21<<12
	return Adrp{
		base: newBase(addr, w),
		rd:   regNameX(w & 0x1f),
		off:  imm21,
		page: uint64(page),
	}
}

func (i Adrp) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("adrp %s, %d ; 0x%x", i.rd, i.off, i.page)
}

func (i Adrp) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, 0x90000000|regBitsX(i.rd)|uint32(i.off&3)<<29|uint32(i.off>>2&0x7ffff)<<5)
}

func (i Adrp) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"adrp",
		i.ObjDump(disasm.DefaultViewCtx()),
		"PC-relative",
		map[string]any{"Rd": i.rd},
	)
}
