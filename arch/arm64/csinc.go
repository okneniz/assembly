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

func NewCsinc(csel Csel) Csinc {
	return Csinc{Csel: csel}
}

func decodeCsinc(w uint32, addr uint64) Instr {
	c, ok := decodeCsel(w, addr).(Csel)
	if !ok {
		// decodeCsel always returns Csel; the branch guards against schema desynchronization
		return Csinc{}
	}

	return Csinc{Csel: c}
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

func (i Csinc) Encode(w io.Writer, pc uint64) (int64, error) {
	return cselWrite(w, i.Csel, csincX, csincW, "csinc")
}

func (i Csinc) MarshalJSON() ([]byte, error) {
	return i.marshal("csinc", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "cond": i.cond})
}
