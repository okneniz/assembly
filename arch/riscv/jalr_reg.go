package riscv

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// JalrReg - the "jalr rs" pseudo-instruction (indirect call): 32-bit jalr ra, 0(rs),
// not compressed (c.jalr decodes back to different text).
type JalrReg struct {
	rs1 string
}

func (i JalrReg) ObjDump(_ disasm.ViewCtx) string {
	return "jalr " + i.rs1
}

func (i JalrReg) Len() int {
	return 4
}

func (i JalrReg) Encode(w io.Writer, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["jalr"][0]|1<<7|regBits(i.rs1)<<15)
}
