package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sub - sub rd, rs1, rs2; pseudo: neg (rs1=zero).
type Sub struct {
	rd, rs1, rs2 string
}

// cSub - compressed forms (c.sub): base - halfword, length 2.
func cSub(rd, rs1, rs2 string) Sub {
	return Sub{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}
}

func (i Sub) ObjDump(_ disasm.ViewCtx) string {
	if i.rs1 == "zero" {
		return fmt.Sprintf("neg %s, %s", i.rd, i.rs2)
	}

	return fmt.Sprintf("sub %s, %s, %s", i.rd, i.rs1, i.rs2)
}

func (i Sub) Encode(w io.Writer, o EncOpts) (int64, error) {
	word := riscvEncodings["sub"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | regBits(i.rs2)<<20
	if half, ok := cR3(i.rd, i.rs1, i.rs2, 0x8C01, 0); ok && !o.NoRVC {
		return writeHalf(w, half)
	}

	return writeWord(w, word)
}

func newSub(ops []Op) (Instr, error) {
	rd, rs1, rs2, err := wantR3(ops, "sub")
	if err != nil {
		return nil, err
	}

	return Sub{
		rd:  rd,
		rs1: rs1,
		rs2: rs2,
	}, nil
}
