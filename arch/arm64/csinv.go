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

func NewCsinv(csel Csel) Csinv {
	return Csinv{Csel: csel}
}

func decodeCsinv(w uint32, addr uint64) Instr {
	c, ok := decodeCsel(w, addr).(Csel)
	if !ok {
		// decodeCsel always returns Csel; the branch guards against schema desynchronization
		return Csinv{}
	}

	return Csinv{Csel: c}
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

func (i Csinv) Encode(w io.Writer, pc uint64) (int64, error) {
	return cselWrite(w, i.Csel, csinvX, csinvW, "csinv")
}

func (i Csinv) MarshalJSON() ([]byte, error) {
	return i.marshal("csinv", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "cond": i.cond})
}
