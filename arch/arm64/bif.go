package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bif — bif.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q).
type Bif struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newBif - the Bif constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newBif(b base, rd, rn, rm VReg, arr string) (Bif, error) {
	err := requireArr("Bif", arr, "8b", "16b")
	if err != nil {
		return Bif{}, err
	}

	return Bif{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const bifEnc uint32 = 786439168 // bif vd, vn, vm (Q=0 form)

func (i Bif) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bif.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Bif) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bif: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("bif: %w", err)
	}

	return writeWord(w, bifEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Bif(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBif(base{}, rd, rn, rm, arr)
}
