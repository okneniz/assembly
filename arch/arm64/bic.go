package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bic — bic.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q).
type Bic struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newBic - the Bic constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newBic(b base, rd, rn, rm VReg, arr string) (Bic, error) {
	err := requireArr("Bic", arr, "8b", "16b")
	if err != nil {
		return Bic{}, err
	}

	return Bic{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const bicEnc uint32 = 241179648 // bic vd, vn, vm (Q=0 form)

func (i Bic) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bic.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Bic) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bic: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("bic: %w", err)
	}

	return writeWord(w, bicEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Bic(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBic(base{}, rd, rn, rm, arr)
}
