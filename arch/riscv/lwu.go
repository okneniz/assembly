package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lwu - lwu rd, off(rs1).
type Lwu struct {
	base

	rd, rs1 string
	off     imm
}

// Lwu - lwu rd, off(rs1).
func (Builder) Lwu(rd, rs1 Reg, off Off) Instr {
	return Lwu{
		rd:  rd.name(),
		rs1: rs1.name(),
		off: immNum(off.v),
	}
}

func decodeLwu(w uint32, addr uint64) Instr {
	return Lwu{
		base: newBase(addr, w),
		rd:   rvRegNames[w>>7&0x1f],
		rs1:  rvRegNames[w>>15&0x1f],
		off:  immNum(iImm(w)),
	}
}

func (i Lwu) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lwu %s, %s(%s)", i.rd, i.off.text(), i.rs1)
}

func (i Lwu) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	off := i.off.val

	bits, err := encI(off)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["lwu"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	return writeWord(w, word)
}

func (i Lwu) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lwu",
		i.ObjDump(disasm.DefaultViewCtx()),
		"RV64I",
		map[string]any{"rd": i.rd, "rs1": i.rs1},
	)
}

func newLwu(ops []Op) (Instr, error) {
	rd, base, off, err := wantR2M(ops, "lwu", false)
	if err != nil {
		return nil, err
	}

	return Lwu{
		rd:  rd,
		rs1: base,
		off: off,
	}, nil
}
