package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srli - srli rd, rs1, shamt; compression: c.srli (shamt < 32).
type Srli struct {
	base

	rd, rs1 string
	shamt   imm
}

// Srli - srli rd, rs1, shamt (shamt5/6 is checked at encoding).
func (Builder) Srli(rd, rs1 Reg, shamt Imm12) Instr {
	return Srli{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

func decodeSrli(w uint32, addr uint64) Instr {
	return Srli{
		base:  newBase(addr, w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt6(w))),
	}
}

// cSrli - compressed forms (c.srli): base - halfword, length 2.
func cSrli(h uint32, addr uint64, rd, rs1 string, shamt int64) Srli {
	return Srli{
		base:  newHalfBase(h, addr),
		rd:    rd,
		rs1:   rs1,
		shamt: immNum(shamt),
	}
}

func (i Srli) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srli %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func (i Srli) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 63 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	word := riscvEncodings["srli"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | uint32(sh)<<20
	if !o.NoRVC {
		// decoder mirror: bit12 - the high bit of shamt; for srli bit12=0 -> sh<32
		if r3ok(i.rd) && i.rd == i.rs1 && sh > 0 && sh < 32 {
			return writeHalf(w, 0x8001|cr3(i.rs1)<<7|uint16(sh&0x1f)<<2) // c.srli
		}
	}

	return writeWord(w, word)
}

func (i Srli) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srli",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newSrli(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "srli")
	if err != nil {
		return nil, err
	}

	return Srli{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
