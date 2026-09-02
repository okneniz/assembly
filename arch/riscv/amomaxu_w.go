package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AmomaxuW - amomaxu.w rd, rs2, (rs1).
type AmomaxuW struct {
	base

	rd, rs1, rs2 string
}

// AmomaxuW - amomaxu.w rd, rs2, (rs1).
func (Builder) AmomaxuW(rd, rs1, rs2 Reg) Instr {
	return AmomaxuW{
		rd:  rd.name(),
		rs1: rs1.name(),
		rs2: rs2.name(),
	}
}

func decodeAmomaxuW(w uint32, addr uint64) Instr {
	return AmomaxuW{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		rs2:  rvRegNames[w>>20&0x1f],
	}
}

func (i AmomaxuW) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("amomaxu.w %s, %s, (%s)", i.rd, i.rs2, i.rs1)
}

func (i AmomaxuW) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["amomaxu_w"][0]|
		regBits(i.rd)<<7|regBits(i.rs1)<<15|regBits(i.rs2)<<20)
}

func (i AmomaxuW) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"amomaxu.w",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32A",
		map[string]any{"rd": i.rd, "rs1": i.rs1, "rs2": i.rs2},
	)
}

func newAmomaxuW(ops []Op) (Instr, error) {
	if len(ops) != 3 {
		return nil, errors.New("amomaxu.w: want rd, rs2, (rs1)")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("amomaxu.w: %w", err)
	}

	rs2, err := wantReg(ops[1], false)
	if err != nil {
		return nil, fmt.Errorf("amomaxu.w: %w", err)
	}

	base, _, err := wantMem(ops[2], "amomaxu.w")
	if err != nil {
		return nil, err
	}

	return AmomaxuW{
		rd:  rd,
		rs1: base,
		rs2: rs2,
	}, nil
}
