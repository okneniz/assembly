package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srai - srai rd, rs1, shamt; compression: c.srai (shamt >= 32).
type Srai struct {
	base

	rd, rs1 string
	shamt   imm
}

// Srai - srai rd, rs1, shamt (shamt5/6 is checked at encoding).
func (Builder) Srai(rd, rs1 Reg, shamt Imm12) Instr {
	return Srai{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

func decodeSrai(w uint32, addr uint64) Instr {
	return Srai{
		base:  newBase(addr, w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt6(w))),
	}
}

// cSrai - compressed forms (c.srai): base - halfword, length 2.
func cSrai(h uint32, addr uint64, rd, rs1 string, shamt int64) Srai {
	return Srai{
		base:  newHalfBase(h, addr),
		rd:    rd,
		rs1:   rs1,
		shamt: immNum(shamt),
	}
}

func (i Srai) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srai %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func (i Srai) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 63 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	word := riscvEncodings["srai"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | uint32(sh)<<20
	if !o.NoRVC {
		// decoder mirror: c.srai is recognized when bit12=1 -> sh >= 32
		if r3ok(i.rd) && i.rd == i.rs1 && sh >= 32 && sh < 64 {
			return writeHalf(w, 0x9001|cr3(i.rs1)<<7|uint16(sh&0x1f)<<2) // c.srai
		}
	}

	return writeWord(w, word)
}

func (i Srai) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srai",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newSrai(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "srai")
	if err != nil {
		return nil, err
	}

	return Srai{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
