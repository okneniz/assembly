package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bltu — bltu rs1, rs2, target.
type Bltu struct {
	base

	rs1, rs2 string
	target   imm // absolute target address
}

func decodeBltu(w uint32, addr uint64) Instr {
	return Bltu{
		base:   newBase(addr, w),
		rs1:    rvRegNames[w>>15&0x1f],
		rs2:    rvRegNames[w>>20&0x1f],
		target: immNum(int64(addr) + bImm(w)),
	}
}

func (i Bltu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bltu %s, %s, %s", i.rs1, i.rs2, i.target.text())
}

func (i Bltu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	target := i.target.val

	bits, err := encB(target - int64(pc))
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["bltu"][0] | regBits(i.rs1)<<15 | regBits(i.rs2)<<20 | bits

	return writeWord(w, word)
}

func (i Bltu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bltu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rs1": i.rs1, "rs2": i.rs2},
	)
}

func newBltu(ops []Op) (Instr, error) {
	rs1, rs2, t, err := wantR2T(ops, "bltu")
	if err != nil {
		return nil, err
	}

	return Bltu{
		rs1:    rs1,
		rs2:    rs2,
		target: t,
	}, nil
}
