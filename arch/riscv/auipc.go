package riscv

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Auipc - auipc rd, imm.
type Auipc struct {
	base

	rd  string
	imm imm
}

// Auipc - auipc rd, imm20 (imm20 = the raw U-type field: rd = pc + imm20<<12).
func (Builder) Auipc(rd Reg, imm Imm20) Instr {
	return Auipc{
		rd:  rd.name(),
		imm: immNum(imm.v),
	}
}

func decodeAuipc(w uint32, addr uint64) Instr {
	return Auipc{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		imm:  immNum(int64(uImm(w))),
	}
}

func (i Auipc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("auipc %s, %s", i.rd, i.imm.text())
}

func (i Auipc) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encU(v)
	if err != nil {
		return 0, err
	}

	return writeWord(w, riscvEncodings["auipc"][0]|regBits(i.rd)<<7|bits)
}

func (i Auipc) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"auipc",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV32I",
		map[string]any{"rd": i.rd},
	)
}

// newAuipc - constructor from parsing: auipc rd, imm.
func newAuipc(ops []Op) (Instr, error) {
	if len(ops) != 2 {
		return nil, errors.New("auipc: want rd, imm")
	}

	rd, err := wantReg(ops[0], false)
	if err != nil {
		return nil, fmt.Errorf("auipc: %w", err)
	}

	e, err := wantExpr(ops[1])
	if err != nil {
		return nil, fmt.Errorf("auipc: %w", err)
	}

	return Auipc{
		rd:  rd,
		imm: e,
	}, nil
}
