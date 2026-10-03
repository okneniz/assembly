package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fneg — fneg fd, fn (double/single by the operand kind).
type Fneg struct {
	rd, rn string
}

// newFneg - the Fneg constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFneg(rd, rn FReg) (Fneg, error) {
	err := requireFpKind("Fneg", rd, rn)
	if err != nil {
		return Fneg{}, err
	}

	return Fneg{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

const (
	fnegD uint32 = 0x1E614000 // fneg dd, dn
	fnegS uint32 = 0x1E214000 // fneg sd, sn
)

func (i Fneg) Encode(w io.Writer) (int64, error) {
	match := fpMatch(i.rd, fnegD, fnegS)

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("fneg: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (i Fneg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fneg %s, %s", i.rd, i.rn)
}
