package riscv

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// JalrReg - the "jalr rs" pseudo-instruction (indirect call): 32-bit jalr ra, 0(rs),
// not compressed (c.jalr decodes back to different text).
type JalrReg struct {
	base

	rs1 string
}

// JalrReg - jalr rs (indirect call: the fixed 32-bit form jalr ra, 0(rs)).
func (Builder) JalrReg(rs1 Reg) Instr {
	return JalrReg{
		rs1: rs1.name(),
	}
}

func (i JalrReg) ObjDump(_ disasm.ViewCtx) string {
	return "jalr " + i.rs1
}

func (i JalrReg) Len() int {
	return 4
}

func (i JalrReg) Encode(w io.Writer, pc uint64, o EncOpts) (int64, error) {
	return writeWord(w, riscvEncodings["jalr"][0]|1<<7|regBits(i.rs1)<<15)
}

func (i JalrReg) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"jalr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Pseudo",
		map[string]any{"rs1": i.rs1},
	)
}
