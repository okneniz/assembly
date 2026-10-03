package loong64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Pcaddu18i - pcaddu18i rd, si20 (1RI20): rd = pc + (si20 << 18). This
// is an address computation, not a jump: the raw si20 is stored.
type Pcaddu18i struct {
	rd  uint8
	imm imm
}

func (i Pcaddu18i) Encode(w io.Writer) (int64, error) {
	word := loongEncodings["pcaddu18i"][0] | uint32(i.rd) | scatterS(i.imm.val, 5, 20)

	return writeWord(w, word)
}

func (i Pcaddu18i) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("pcaddu18i %s, %s", laRegName(i.rd), i.imm.text())
}
