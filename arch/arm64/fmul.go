package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fmul — fmul fd, fn, fm (double/single by the operand kind).
type Fmul struct {
	rd, rn, rm string
}

// newFmul - the Fmul constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFmul(rd, rn, rm FReg) (Fmul, error) {
	err := requireFpKind("Fmul", rd, rn, rm)
	if err != nil {
		return Fmul{}, err
	}

	return Fmul{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const (
	fmulD uint32 = 0x1E600800 // fmul dd, dn, dm
	fmulS uint32 = 0x1E200800 // fmul sd, sn, sm
)

func (i Fmul) Encode(w io.Writer) (int64, error) {
	match := fpMatch(i.rd, fmulD, fmulS)

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fmul: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Fmul) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fmul %s, %s, %s", i.rd, i.rn, i.rm)
}
