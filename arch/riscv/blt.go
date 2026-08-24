package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Blt - blt rs1, rs2, target; pseudo: bgtz (rs1=zero), bltz (rs2=zero).
type Blt struct {
	base

	rs1, rs2 string
	target   imm // absolute target address
}

func decodeBlt(w uint32, addr uint64) Instr {
	return Blt{
		base:   newBase(addr, w),
		rs1:    rvRegNames[w>>15&0x1f],
		rs2:    rvRegNames[w>>20&0x1f],
		target: immNum(int64(addr) + bImm(w)),
	}
}

func (i Blt) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("bgtz %s, %s", i.rs2, i.target.text())
	}

	if i.rs2 == "zero" {
		return fmt.Sprintf("bltz %s, %s", i.rs1, i.target.text())
	}

	return fmt.Sprintf("blt %s, %s, %s", i.rs1, i.rs2, i.target.text())
}

func (i Blt) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	target := i.target.val

	bits, err := encB(target - int64(pc))
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["blt"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func (i Blt) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"blt",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newBlt(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "blt")
	if err != nil {
		return nil, err
	}

	return Blt{
		rs1:    rs1,
		rs2:    rs2,
		target: t,
	}, nil
}
