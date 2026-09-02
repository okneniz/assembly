package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

const (
	csnegX uint32 = 0xDA800400
	csnegW uint32 = 0x5A800400
)

// Csneg — csneg rd, rn, rm, cond; pseudo: cneg (rn == rm).
type Csneg struct {
	Csel
}

// Csneg — csneg rd, rn, rm, cond (cneg pseudo); the operand
// constraints are those of Csel.
func (Builder) Csneg(rd, rn, rm Reg, cond string) (Instr, error) {
	base, err := New().Csel(rd, rn, rm, cond)
	if err != nil {
		return nil, err
	}

	c, ok := base.(Csel)
	if !ok {
		return nil, fmt.Errorf("csneg: internal: want Csel, got %T", base)
	}

	return Csneg{Csel: c}, nil
}

func decodeCsneg(w uint32, addr uint64) Instr {
	c, ok := decodeCsel(w, addr).(Csel)
	if !ok {
		// decodeCsel always returns Csel; the branch guards against schema desynchronization
		return Csneg{}
	}

	return Csneg{Csel: c}
}

func (i Csneg) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == i.rm {
		return fmt.Sprintf("cneg %s, %s, %s", i.rd, i.rm, invertCond(i.cond))
	}

	return fmt.Sprintf("csneg %s, %s, %s, %s", i.rd, i.rn, i.rm, i.cond)
}

func (i Csneg) Encode(w io.Writer, pc uint64) (int64, error) {
	return cselWrite(w, i.Csel, csnegX, csnegW, "csneg")
}

func (i Csneg) MarshalJSON() ([]byte, error) {
	return i.marshal("csneg", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "cond": i.cond})
}
