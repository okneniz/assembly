package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Eor — eor.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q).
type Eor struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newEor - the Eor constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newEor(b base, rd, rn, rm VReg, arr string) (Eor, error) {
	err := requireArr("Eor", arr, "8b", "16b")
	if err != nil {
		return Eor{}, err
	}

	return Eor{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const eorEnc uint32 = 773856256 // eor vd, vn, vm (Q=0 form)

func (i Eor) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("eor.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Eor) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("eor: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("eor: %w", err)
	}

	return writeWord(w, eorEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Eor(rd, rn, rm VReg, arr string) (Instr, error) {
	return newEor(base{}, rd, rn, rm, arr)
}
