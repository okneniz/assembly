package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bl - bl offs (1RI26, the immediate split d10k16): rd = $ra implicitly
// ($ra = pc + 4), jump to pc + offs (word-scaled). The decoded form
// stores the byte offset itself.
type Bl struct {
	base

	off imm
}

// NewBl - bl offs (the pc-relative byte offset).
func NewBl(off int64) Instr {
	return Bl{
		off: immNum(off),
	}
}

func decodeBl(w uint32, addr uint64) Instr {
	return Bl{
		base: newBase(addr, w),
		off:  immNum(d10k16Imm(w) << 2),
	}
}

func (i Bl) ObjDump(_ disasm.ViewCtx) string {
	return "bl " + i.off.text()
}

func (i Bl) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 26, "bl offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["bl"][0] | scatterD10k16(off)

	return writeWord(w, word)
}

func (i Bl) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"bl",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"off": i.off.val},
	)
}
