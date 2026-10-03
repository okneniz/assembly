package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Aesmc — aesmc.16b vd, vn (the AES round helper; .16b implicit).
type Aesmc struct {
	rd, rn string
}

// newAesmc - the Aesmc constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newAesmc(rd, rn VReg) Aesmc {
	return Aesmc{
		rd: rd.name(),
		rn: rn.name(),
	}
}

const aesmcEnc uint32 = 1311270912 // aesmc vd, vn

func (i Aesmc) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("aesmc: %w", err)
	}

	return writeWord(w, aesmcEnc|rd|rn<<5)
}

func (i Aesmc) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("aesmc.16b %s, %s", i.rd, i.rn)
}
