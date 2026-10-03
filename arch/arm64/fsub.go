package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Fsub — fsub fd, fn, fm (double/single by the operand kind).
type Fsub struct {
	rd, rn, rm string
}

// newFsub - the Fsub constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newFsub(rd, rn, rm FReg) (Fsub, error) {
	err := requireFpKind("Fsub", rd, rn, rm)
	if err != nil {
		return Fsub{}, err
	}

	return Fsub{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

const (
	fsubD uint32 = 0x1E603800 // fsub dd, dn, dm
	fsubS uint32 = 0x1E203800 // fsub sd, sn, sm
)

func (i Fsub) Encode(w io.Writer) (int64, error) {
	match := fpMatch(i.rd, fsubD, fsubS)

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("fsub: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Fsub) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("fsub %s, %s, %s", i.rd, i.rn, i.rm)
}
