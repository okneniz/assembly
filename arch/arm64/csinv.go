package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

const (
	csinvX uint32 = 0xDA800000
	csinvW uint32 = 0x5A800000
)

// Csinv — csinv rd, rn, rm, cond; pseudo: csetm/cinv (inverted cond).
type Csinv struct {
	Csel
}

func (i Csinv) ObjDump(_ disasm.ViewCtx) string {
	zr := zeroReg(i.rd)
	inv := invertCond(i.cond)
	if i.rn == zr && i.rm == zr {
		return fmt.Sprintf("csetm %s, %s", i.rd, inv)
	}

	if i.rn == i.rm {
		return fmt.Sprintf("cinv %s, %s, %s", i.rd, i.rm, inv)
	}

	return fmt.Sprintf("csinv %s, %s, %s, %s", i.rd, i.rn, i.rm, i.cond)
}

func (i Csinv) Encode(w io.Writer) (int64, error) {
	return cselWrite(w, i.Csel, csinvX, csinvW, "csinv")
}
