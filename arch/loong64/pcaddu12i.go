package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Pcaddu12i - pcaddu12i rd, si20 (1RI20): rd = pc + (si20 << 12). This
// is an address computation, not a jump: the raw si20 is stored.
type Pcaddu12i struct {
	base

	rd  uint8
	imm imm
}

// Pcaddu12i - pcaddu12i rd, si20.
func (Builder) Pcaddu12i(rd Reg, v Imm20) Instr {
	return Pcaddu12i{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

func decodePcaddu12i(w uint32, addr uint64) Instr {
	return Pcaddu12i{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func (i Pcaddu12i) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("pcaddu12i %s, %s", laRegName(i.rd), i.imm.text())
}

func (i Pcaddu12i) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["pcaddu12i"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Pcaddu12i) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"pcaddu12i",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "imm": i.imm.val},
	)
}
