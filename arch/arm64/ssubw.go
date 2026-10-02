package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ssubw — ssubw{,2}.Arr vd, vn, vm (the widening three-same group:
// the arrangement is the RESULT's - one lane wider than the Rm source;
// Q=1 prints the "2" form).
type Ssubw struct {
	q, size    uint32
	rd, rn, rm string
}

// newSsubw - the Ssubw constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSsubw(q, size uint32, rd, rn, rm VReg) (Ssubw, error) {
	if size > 2 {
		return Ssubw{}, fmt.Errorf(
			"arm64.NewSsubw: the source arrangement is too narrow for .2d results",
		)
	}

	return Ssubw{
		q:    q,
		size: size,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
	}, nil
}

const ssubwEnc uint32 = 236990464 // ssubw vd, vn, vm (Q=0 form)

func (i Ssubw) ObjDump(_ disasm.ViewCtx) string {
	name := "ssubw"
	if i.q == 1 {
		name += "2"
	}

	return fmt.Sprintf("%s.%s %s, %s, %s",
		name, decodeArrangement(i.q, i.size+1), i.rd, i.rn, i.rm)
}

func (i Ssubw) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("ssubw: %w", err)
	}

	return writeWord(w, ssubwEnc|i.q<<30|i.size<<22|rd|rn<<5|rm<<16)
}
