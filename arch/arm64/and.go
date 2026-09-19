package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// And — and.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q; Rn == Rm prints as the mov alias).
type And struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newAnd - the And constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newAnd(b base, rd, rn, rm VReg, arr string) (And, error) {
	err := requireArr("And", arr, "8b", "16b")
	if err != nil {
		return And{}, err
	}

	return And{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const andEnc uint32 = 236985344 // and vd, vn, vm (Q=0 form)

func (i And) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == i.rm { // the mov alias: and/orr with Rn == Rm
		return fmt.Sprintf("mov.%s %s, %s", i.arr, i.rd, i.rm)
	}

	return fmt.Sprintf("and.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i And) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	return writeWord(w, andEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) And(rd, rn, rm VReg, arr string) (Instr, error) {
	return newAnd(base{}, rd, rn, rm, arr)
}
