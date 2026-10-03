package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Addp — addp.Arr vd, vn, vm (the three-same group: the lane width
// rides the arrangement).
type Addp struct {
	rd, rn, rm string
	arr        string // 8b/16b/4h/8h/2s/4s/2d
}

// newAddp - the Addp constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newAddp(rd, rn, rm VReg, arr string) (Addp, error) {
	err := requireArr("Addp", arr, "8b", "16b", "4h", "8h", "2s", "4s", "2d")
	if err != nil {
		return Addp{}, err
	}

	return Addp{
		rd:  rd.name(),
		rn:  rn.name(),
		rm:  rm.name(),
		arr: arr,
	}, nil
}

const addpEnc uint32 = 237026304 // addp vd, vn, vm (Q=0, size=0 form)

func (i Addp) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("addp: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("addp: %w", err)
	}

	return writeWord(w, addpEnc|q<<30|size<<22|rd|rn<<5|rm<<16)
}

func (i Addp) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("addp.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}
