package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Bit — bit.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q).
type Bit struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newBit - the Bit constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newBit(b base, rd, rn, rm VReg, arr string) (Bit, error) {
	err := requireArr("Bit", arr, "8b", "16b")
	if err != nil {
		return Bit{}, err
	}

	return Bit{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const bitEnc uint32 = 782244864 // bit vd, vn, vm (Q=0 form)

func (i Bit) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("bit.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Bit) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("bit: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("bit: %w", err)
	}

	return writeWord(w, bitEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Bit(rd, rn, rm VReg, arr string) (Instr, error) {
	return newBit(base{}, rd, rn, rm, arr)
}
