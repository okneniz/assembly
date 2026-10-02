package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

const (
	csincX uint32 = 0x9A800400
	csincW uint32 = 0x1A800400
)

// Csinc — csinc rd, rn, rm, cond; pseudo: cset rd, cond (rn=rm=zr),
// cinc rd, rm, cond (rn == rm) — with an inverted condition.
type Csinc struct {
	Csel
}

func (i Csinc) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	inv := invertCond(i.cond)
	if i.rn == zr && i.rm == zr {
		return fmt.Sprintf("cset %s, %s", i.rd, inv)
	}

	if i.rn == i.rm {
		return fmt.Sprintf("cinc %s, %s, %s", i.rd, i.rm, inv)
	}

	return fmt.Sprintf("csinc %s, %s, %s, %s", i.rd, i.rn, i.rm, i.cond)
}

func (i Csinc) Encode(w io.Writer) (int64, error) {
	return cselWrite(w, i.Csel, csincX, csincW, "csinc")
}
