package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Orn — orn.8b/16b vd, vn, vm (the logical three-same group:
// bits 23:22 are the opcode, the arrangement is 8b/16b by Q).
type Orn struct {
	base

	rd, rn, rm string
	arr        string // 8b/16b
}

// newOrn - the Orn constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newOrn(b base, rd, rn, rm VReg, arr string) (Orn, error) {
	err := requireArr("Orn", arr, "8b", "16b")
	if err != nil {
		return Orn{}, err
	}

	return Orn{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		rm:   rm.name(),
		arr:  arr,
	}, nil
}

const ornEnc uint32 = 249568256 // orn vd, vn, vm (Q=0 form)

func (i Orn) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("orn.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}

func (i Orn) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("orn: %w", err)
	}

	q, _, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("orn: %w", err)
	}

	return writeWord(w, ornEnc|q<<30|rd|rn<<5|rm<<16)
}

func (Builder) Orn(rd, rn, rm VReg, arr string) (Instr, error) {
	return newOrn(base{}, rd, rn, rm, arr)
}
