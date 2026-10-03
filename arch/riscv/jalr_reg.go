package riscv

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// JalrReg - the "jalr rs" pseudo-instruction (indirect call): jalr
// ra, 0(rs); compresses to c.jalr when rs1 is not x0 (the decoder
// prints c.jalr back as "jalr rs" through the Jalr canon).
type JalrReg struct {
	rs1 string
}

func (i JalrReg) Encode(w io.Writer, o EncOpts) (int64, error) {
	if !o.NoRVC {
		if r := r5(i.rs1); r != 0 {
			return writeHalf(w, 0x9002|r<<7) // c.jalr (llvm-mc parity)
		}
	}

	return writeWord(w, riscvEncodings["jalr"][0]|1<<7|regBits(i.rs1)<<15)
}

func (i JalrReg) Len() int {
	if r := r5(i.rs1); r != 0 {
		return 2 // c.jalr
	}

	return 4
}

func (i JalrReg) ObjDump(_ disasm.ViewCtx) string {
	return "jalr " + i.rs1
}
