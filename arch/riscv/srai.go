package riscv

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Srai - srai rd, rs1, shamt; compression: c.srai (shamt >= 32).
type Srai struct {
	rd, rs1 string
	shamt   imm
}

// cSrai - compressed forms (c.srai): base - halfword, length 2.
func cSrai(rd, rs1 string, shamt int64) Srai {
	return Srai{
		rd:    rd,
		rs1:   rs1,
		shamt: immNum(shamt),
	}
}

func (i Srai) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("srai %s, %s, %s", i.rd, i.rs1, i.shamt.text())
}

func (i Srai) Encode(w io.Writer, o EncOpts) (int64, error) {
	sh := i.shamt.val

	if sh < 0 || sh > 63 {
		return 0, fmt.Errorf("shift amount %d out of range", sh)
	}

	word := riscvEncodings["srai"][0] | regBits(i.rd)<<7 | regBits(i.rs1)<<15 | uint32(sh)<<20
	if !o.NoRVC {
		// decoder mirror: c.srai is recognized when bit12=1 -> sh >= 32;
		// the CI fixed bits are b10=1, b6=0 (the 0x8401 class, llvm-mc
		// parity: srai a2, 55 = 0x965d)
		if r3ok(i.rd) && i.rd == i.rs1 && sh >= 32 && sh < 64 {
			return writeHalf(
				w,
				0x8401|uint16(sh&0x20)<<7|cr3(i.rs1)<<7|uint16(sh&0x1f)<<2,
			) // c.srai
		}
	}

	return writeWord(w, word)
}

func newSrai(ops []Op) (Instr, error) {
	rd, rs1, m, err := wantI3(ops, "srai")
	if err != nil {
		return nil, err
	}

	return Srai{
		rd:    rd,
		rs1:   rs1,
		shamt: m,
	}, nil
}
