package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

const (
	msubX uint32 = 0x9B008000
	msubW uint32 = 0x1B008000
)

// Msub — msub rd, rn, rm, ra; pseudo: mneg (ra = xzr).
type Msub struct {
	Madd
}

func (i Msub) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, msubX, msubW)
	if err != nil {
		return 0, fmt.Errorf("msub: %w", err)
	}

	return msubWrite(w, match, i.Madd)
}

func (i Msub) ObjDump(_ disasm.ViewCtx) string {
	zr := "xzr"
	if i.rd[0] == 'w' {
		zr = "wzr"
	}

	if i.ra == zr {
		return fmt.Sprintf("mneg %s, %s, %s", i.rd, i.rn, i.rm)
	}

	return fmt.Sprintf("msub %s, %s, %s, %s", i.rd, i.rn, i.rm, i.ra)
}
