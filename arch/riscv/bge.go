package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bge - bge rs1, rs2, target; pseudo: blez (rs1=zero), bgez (rs2=zero).
type Bge struct {
	base

	rs1, rs2 string
	target   imm // absolute target address
}

func decodeBge(w uint32, addr uint64) Instr {
	return Bge{
		base:   newBase(addr, w),
		rs1:    rvRegNames[w>>15&0x1f],
		rs2:    rvRegNames[w>>20&0x1f],
		target: immNum(int64(addr) + bImm(w)),
	}
}

func (i Bge) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("blez %s, %s", i.rs2, i.target.text())
	}

	if i.rs2 == "zero" {
		return fmt.Sprintf("bgez %s, %s", i.rs1, i.target.text())
	}

	return fmt.Sprintf("bge %s, %s, %s", i.rs1, i.rs2, i.target.text())
}

func (i Bge) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	target := i.target.val

	bits, err := encB(target - int64(pc))
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bge"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func (i Bge) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bge",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newBge(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bge")
	if err != nil {
		return nil, err
	}

	return Bge{
		rs1:    rs1,
		rs2:    rs2,
		target: t,
	}, nil
}
