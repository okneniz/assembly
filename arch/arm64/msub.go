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

// Msub - msub rd, rn, rm, ra; pseudo: mneg (ra = xzr).
type Msub struct {
	Madd
}

func NewMsub(madd Madd) Msub {
	return Msub{Madd: madd}
}

func decodeMsub(w uint32, addr uint64) Instr {
	m, ok := decodeMadd(w, addr).(Madd)
	if !ok {
		// decodeMadd always returns Madd; this branch guards against schema desynchronization
		return Msub{}
	}

	return Msub{Madd: m}
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

func (i Msub) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, msubX, msubW)
	if err != nil {
		return 0, fmt.Errorf("msub: %w", err)
	}

	return msubWrite(w, match, i.Madd)
}

func (i Msub) MarshalJSON() ([]byte, error) {
	return i.marshal("msub", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "Ra": i.ra})
}
