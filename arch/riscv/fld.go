package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fld - fld rd, off(rs1) (rd - floating).
type Fld struct {
	base

	rd, rs1 string
	off     imm
}

func decodeFld(w uint32, addr uint64) Instr {
	return Fld{
		base: newBase(addr, w),
		rd:   rvFRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Fld) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fld %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Fld) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["fld"][0] | fregBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Fld) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"fld",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64D",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newFld(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "fld", true)
	if err != nil {
		return nil, err
	}

	return Fld{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
