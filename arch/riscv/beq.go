package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Beq - beq rs1, rs2, target; pseudo: beqz.
type Beq struct {
	base

	rs1, rs2 string
	target   imm // absolute target address
}

// Beq - beq rs1, rs2, target (absolute target address).
func (Builder) Beq(rs1, rs2 Reg, target int64) Instr {
	return Beq{
		rs1:    rs1.name(),
		rs2:    rs2.name(),
		target: immNum(target),
	}
}

func decodeBeq(w uint32, addr uint64) Instr {
	return Beq{
		base:   newBase(addr, w),
		rs1:    rvRegNames[w>>15&0x1f],
		rs2:    rvRegNames[w>>20&0x1f],
		target: immNum(int64(addr) + bImm(w)),
	}
}

// cBeq - compressed forms (c.beqz): base - halfword, length 2.
func cBeq(h uint32, addr uint64, rs1, rs2 string, target int64) Beq {
	return Beq{
		base:   newHalfBase(h, addr),
		rs1:    rs1,
		rs2:    rs2,
		target: immNum(target),
	}
}

func (i Beq) ObjDump(_ disasm.ViewCtx) string {
	if i.rs2 == "zero" {
		return fmt.Sprintf("beqz %s, %s", i.rs1, i.target.text())
	}

	return fmt.Sprintf("beq %s, %s, %s", i.rs1, i.rs2, i.target.text())
}

func (i Beq) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	target := i.target.val

	bits, err := encB(target - int64(pc))
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["beq"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	if !o.NoRVC {
		if half, ok := cbeqz(i.rs1, i.rs2, "beq", target-int64(pc)); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

func (i Beq) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"beq",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newBeq(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "beq")
	if err != nil {
		return nil, err
	}

	return Beq{
		rs1:    rs1,
		rs2:    rs2,
		target: t,
	}, nil
}
