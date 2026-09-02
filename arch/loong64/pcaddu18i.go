package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Pcaddu18i - pcaddu18i rd, si20 (1RI20): rd = pc + (si20 << 18). This
// is an address computation, not a jump: the raw si20 is stored.
type Pcaddu18i struct {
	base

	rd  uint8
	imm imm
}

// Pcaddu18i - pcaddu18i rd, si20.
func (Builder) Pcaddu18i(rd Reg, v Imm20) Instr {
	return Pcaddu18i{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

func decodePcaddu18i(w uint32, addr uint64) Instr {
	return Pcaddu18i{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func (i Pcaddu18i) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("pcaddu18i %s, %s", laRegName(i.rd), i.imm.text())
}

func (i Pcaddu18i) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["pcaddu18i"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Pcaddu18i) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"pcaddu18i",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "imm": i.imm.val},
	)
}
