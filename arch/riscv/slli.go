package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Slli - slli rd, rs1, shamt; compression: c.slli.
type Slli struct {
	rd, rs1 string
	shamt   imm
}

// cSlli - compressed forms (c.slli): base - halfword, length 2.
func cSlli(rd, rs1 string, shamt int64) Slli {
	return Slli{
		rd:    rd,
		rs1:   rs1,
		shamt: immNum(shamt),
	}
}

func (i Slli) Encode(w io.Writer, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 63 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	word := riscvEncodings["slli"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | uint32(sh)<<20
	if !o.NoRVC {
		if r := r5(i.rd); i.rd == i.rs1 && r != 0 && sh > 0 && sh < 64 {
			u := uint16(sh)
			return writeHalf(w, 0x0002|r<<7|(u>>5&1)<<12|(u&0x1f)<<2) // c.slli
		}
	}

	return writeWord(w, word)
}

func (i Slli) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("slli %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func newSlli(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "slli")
	if err != nil {
		return nil, err
	}

	return Slli{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
