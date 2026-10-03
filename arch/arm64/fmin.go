package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fmin — fmin fd, fn, fm (double/single by the operand kind).
type Fmin struct {
	rd, rn, rm string
}

// newFmin - the Fmin constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmin(rd, rn, rm FReg) (Fmin, error) {
	err := requireFpKind("Fmin", rd, rn, rm)
	if err != nil {
		return Fmin{}, err
	}

	return Fmin{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const (
	fminD uint32 = 0x1E605800 // fmin dd, dn, dm
	fminS uint32 = 0x1E205800 // fmin sd, sn, sm
)

func (i Fmin) Encode(w io.Writer) (int64, error) {
	match := fpMatch(i.rd, fminD, fminS)

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmin: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Fmin) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmin %s, %s, %s", i.rd, i.rn, i.rm)
}
