package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Addiw - addiw rd, rs1, imm; pseudo: sext.w (imm=0); compression: c.addiw.
type Addiw struct {
	rd, rs1 string
	imm     imm
}

// cAddiw - compressed forms (c.addiw): base - halfword, length 2.
func cAddiw(rd, rs1 string, imm int64) Addiw {
	return Addiw{
		rd:  rd,
		rs1: rs1,
		imm: immNum(imm),
	}
}

func (i Addiw) ObjDump(_ disasm.ViewCtx) string {
	if i.imm.val == 0 {
		return fmt.Sprintf("sext.w %s, %s", i.rd, i.rs1)
	}

	return fmt.Sprintf("addiw %s, %s, %s", i.rd, i.rs1, i.imm.text())
}

func (i Addiw) Encode(w io.Writer, o EncOpts) (int64, error) {
	v := i.imm.val

	bits, err := encI(v)
	if err != nil {
		return 0, err
	}

	word := riscvEncodings["addiw"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | bits
	if !o.NoRVC {
		if r := r5(i.rd); i.rd == i.rs1 && r != 0 && fits6(v) {
			return writeHalf(w, 0x2001|r<<7|ciBits(v)) // c.addiw
		}
	}

	return writeWord(w, word)
}

func newAddiw(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "addiw")
	if err != nil {
		return nil, err
	}

	return Addiw{
		rd:  rd,
		rs1: rs1,
		imm: m,
	}, nil
}
