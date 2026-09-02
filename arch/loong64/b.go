package loong64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// B - b offs (1RI26, the immediate split d10k16): unconditional jump to
// pc + offs (word-scaled). The off field holds the byte offset itself.
type B struct {
	base

	off imm
}

// B - b offs (the pc-relative byte offset).
func (Builder) B(off int64) Instr {
	return B{
		off: immNum(off),
	}
}

func decodeB(w uint32, addr uint64) Instr {
	return B{
		base: newBase(addr, w),
		off:  immNum(d10k16Imm(w) << 2),
	}
}

func (i B) ObjDump(_ disasm.ViewCtx) string {
	return "b " + i.off.text()
}

func (i B) Encode(w io.Writer, pc uint64) (int64, error) {
	off, err := encPs2(i.off.val, 26, "b offset")
	if err != nil {
		return 0, err
	}

	word := loongEncodings["b"][0] | scatterD10k16(off)

	return writeWord(w, word)
}

func (i B) MarshalJSON() ([]byte, error) {
	return i.marshalDTO(
		"b",
		i.ObjDump(disasm.DefaultViewCtx()),
		"LA64",
		map[string]any{"off": i.off.val},
	)
}
