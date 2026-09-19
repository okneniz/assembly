package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Orr — orr.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q; Rn == Rm prints as the mov alias).
type Orr struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newOrr - the Orr constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newOrr(b base, rd, rn, rm VReg, arr string) (Orr, error) {
	err := requireArr("Orr", arr, "8b", "16b")
	if err != nil {
		return Orr{}, err
	}

	return Orr{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const orrEnc uint32 = 245373952 // orr vd, vn, vm (Q=0 form)

func (i Orr) ObjDump(_ disasm.ViewCtx) string {
	if i.rn == i.rm { // the mov alias: and/orr with Rn == Rm
		return fmt.Sprintf("mov.%s %s, %s", i.arr, i.rd, i.rm)
	}

	return fmt.Sprintf("orr.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Orr) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	return writeWord(w, orrEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Orr(rd, rn, rm VReg, arr string) (Instr, error) {
	return newOrr(base{}, rd, rn, rm, arr)
}
