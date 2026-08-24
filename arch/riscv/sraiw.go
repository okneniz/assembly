package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sraiw — sraiw rd, rs1, shamt (shamt5).
type Sraiw struct {
	base

	rd, rs1 string
	shamt   imm
}

func decodeSraiw(w uint32, addr uint64) Instr {
	return Sraiw{
		base:  newBase(addr, w),
		rd:    rvRegNames[w>>7&0x1f],
		rs1:   rvRegNames[w>>15&0x1f],
		shamt: immNum(int64(shamt5(w))),
	}
}

func (i Sraiw) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sraiw %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func (i Sraiw) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 31 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	return writeWord(
		w,
		riscvEncodings["sraiw"][0]|regBits(i.rd)<<7|regBits(i.rs1)<<15|uint32(sh)<<20,
	)
}

func (i Sraiw) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"sraiw",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newSraiw(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "sraiw")
	if err != nil {
		return nil, err
	}

	return Sraiw{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
