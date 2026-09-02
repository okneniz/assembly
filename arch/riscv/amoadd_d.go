package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmoaddD - amoadd.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] += rs2.
type AmoaddD struct {
	base

	rd, rs1, rs2 string
}

// AmoaddD - amoadd.d rd, rs2, (rs1): rd = old MEM[rs1]; MEM[rs1] += rs2.
func (Builder) AmoaddD(rd, rs1, rs2 Reg) Instr {
	return AmoaddD{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeAmoaddD(w uint32, addr uint64) Instr {
	return AmoaddD{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i AmoaddD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amoadd.d %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmoaddD) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amoadd_d"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func (i AmoaddD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amoadd.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64A",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAmoaddD(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amoadd.d: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amoadd.d: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amoadd.d: %w", err)
	}

	base, _, err := wantMem(ops[2], "amoadd.d")
	if err != nil {
		return nil, err
	}

	return AmoaddD{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
