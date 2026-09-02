package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Lu32iD - lu32i.d rd, si20 (1RI20): rd = si20 << 32 concatenated into
// bits 51:32.
type Lu32iD struct {
	base

	rd  uint8
	imm imm
}

// Lu32iD - lu32i.d rd, si20.
func (Builder) Lu32iD(rd Reg, v Imm20) Instr {
	return Lu32iD{
		rd:  rd.Num(),
		imm: immNum(v.Val()),
	}
}

func decodeLu32iD(w uint32, addr uint64) Instr {
	return Lu32iD{
		base: newBase(addr, w),
		rd:   uint8(w & 0x1f),
		imm:  immNum(sField(w, 5, 20)),
	}
}

func (i Lu32iD) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lu32i.d %s, %s", laRegName(i.rd), i.imm.text())
}

func (i Lu32iD) Encode(w io.Writer, _ uint64) (int64, error) {
	word := loongEncodings["lu32i.d"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Lu32iD) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"lu32i.d",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"rd": laRegName(i.rd), "imm": i.imm.val},
	)
}
