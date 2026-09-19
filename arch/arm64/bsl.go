package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bsl — bsl.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q).
type Bsl struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newBsl - the Bsl constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newBsl(b base, rd, rn, rm VReg, arr string) (Bsl, error) {
	err := requireArr("Bsl", arr, "8b", "16b")
	if err != nil {
		return Bsl{}, err
	}

	return Bsl{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const bslEnc uint32 = 778050560 // bsl vd, vn, vm (Q=0 form)

func (i Bsl) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bsl.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Bsl) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bsl: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("bsl: %w", err)
	}

	return writeWord(w, bslEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Bsl(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBsl(base{}, rd, rn, rm, arr)
}
