package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Usubw — usubw{,2}.Arr vd, vn, vm (the widening three-same group:
// the arrangement is the RESULT's - one lane wider than the Rm source;
// Q=1 prints the "2" form).
type Usubw struct {
	q, size    uint32
	rd, rn, rm string
}

// newUsubw - the Usubw constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newUsubw(q, size uint32, rd, rn, rm VReg) (Usubw, error) {
	if size > 2 {
		return Usubw{}, errors.New(
			"arm64.NewUsubw: the source arrangement is too narrow for .2d results",
		)
	}

	return Usubw{
		q:    q,
		size: size,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const usubwEnc uint32 = 773861376 // usubw vd, vn, vm (Q=0 form)

func (i Usubw) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("usubw: %w", err)
	}

	return writeWord(w, usubwEnc|i.q<<30|i.size<<22|rd|rn<<5|rm<<16)
}

func (i Usubw) ObjDump(_ disasm.ViewCtx) string {
	name := "usubw"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s",
		name, decodeArrangement(i.q, i.size+1), i.rd, i.rn, i.rm)
}
