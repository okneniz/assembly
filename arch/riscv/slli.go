package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slli - slli rd, rs1, shamt; compression: c.slli.
type Slli struct {
	base

	rd, rs1 string
	shamt   imm
}

// Slli - slli rd, rs1, shamt (shamt5/6 is checked at encoding).
func (Builder) Slli(rd, rs1 Reg, shamt Imm12) Instr {
	return Slli{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

func decodeSlli(w uint32, addr uint64) Instr {
	return Slli{
		base:  newBase(addr, w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt6(w))),
	}
}

// cSlli - compressed forms (c.slli): base - halfword, length 2.
func cSlli(h uint32, addr uint64, rd, rs1 string, shamt int64) Slli {
	return Slli{
		base:  newHalfBase(h, addr),
		rd:    rd,
		rs1:   rs1,
		shamt: immNum(shamt),
	}
}

func (i Slli) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slli %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func (i Slli) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 63 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	word := riscvEncodings["slli"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | uint32(sh)<<20
	if !o.NoRVC {
		if r := r5(i.rd); i.rd == i.rs1 && r != 0 && sh > 0 && sh < 64 {
			u := uint16(sh)
			return writeHalf(w, 0x0002|r<<7|(u>>5&1)<<12|(u&0x1f)<<2) // c.slli
		}
	}

	return writeWord(w, word)
}

func (i Slli) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"slli",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newSlli(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "slli")
	if err != nil {
		return nil, err
	}

	return Slli{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
