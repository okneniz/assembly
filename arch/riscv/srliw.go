package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srliw - srliw rd, rs1, shamt (shamt5).
type Srliw struct {
	base

	rd, rs1 string
	shamt   imm
}

// Srliw - srliw rd, rs1, shamt (shamt5).
func (Builder) Srliw(rd, rs1 Reg, shamt Imm12) Instr {
	return Srliw{
		rd:    rd.name(),
		rs1:   rs1.name(),
		shamt: immNum(shamt.v),
	}
}

func decodeSrliw(w uint32, addr uint64) Instr {
	return Srliw{
		base:  newBase(addr, w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt5(w))),
	}
}

func (i Srliw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srliw %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func (i Srliw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 31 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	return writeWord(
		w,
		riscvEncodings["srliw"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|uint32(sh)<<20,
	)
}

func (i Srliw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"srliw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newSrliw(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "srliw")
	if err != nil {
		return nil, err
	}

	return Srliw{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
