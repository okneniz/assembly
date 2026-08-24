package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bne - bne rs1, rs2, target; pseudo: bnez.
type Bne struct {
	base

	rs1, rs2 string
	target   imm // absolute target address
}

func decodeBne(w uint32, addr uint64) Instr {
	return Bne{
		base:   newBase(addr, w),
		rs1:    rvRegNames[w>>15&0x1f],
		rs2:    rvRegNames[w>>20&0x1f],
		target: immNum(int64(addr) + bImm(w)),
	}
}

// cBne - compressed forms (c.bnez): base - halfword, length 2.
func cBne(h uint32, addr uint64, rs1, rs2 string, target int64) Bne {
	return Bne{
		base:   newHalfBase(h, addr),
		rs1:    rs1,
		rs2:    rs2,
		target: immNum(target),
	}
}

func (i Bne) ObjDump(_ disasm.ViewCtx) string {
	if i.rs2 == "zero" {
		return fmt.Sprintf("bnez %s, %s", i.rs1, i.target.text())
	}

	return fmt.Sprintf("bne %s, %s, %s", i.rs1, i.rs2, i.target.text())
}

func (i Bne) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	target := i.target.val

	bits, err := encB(target - int64(pc))
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bne"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits
	if !o.NoRVC {
		if half, ok := cbeqz(i.rs1, i.rs2, "bne", target-int64(pc)); ok {
			return writeHalf(w, half)
		}
	}

	return writeWord(w, word)
}

func (i Bne) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bne",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newBne(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bne")
	if err != nil {
		return nil, err
	}

	return Bne{
		rs1:    rs1,
		rs2:    rs2,
		target: t,
	}, nil
}
