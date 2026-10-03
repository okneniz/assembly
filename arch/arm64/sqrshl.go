package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sqrshl — sqrshl.Arr vd, vn, vm (the three-same group: the lane width
// rides the arrangement).
type Sqrshl struct {
	rd, rn, rm string
	arr        string // 8b/16b/4h/8h/2s/4s/2d
}

// newSqrshl - the Sqrshl constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newSqrshl(rd, rn, rm VReg, arr string) (Sqrshl, error) {
	err := requireArr("Sqrshl", arr, "8b", "16b", "4h", "8h", "2s", "4s", "2d")
	if err != nil {
		return Sqrshl{}, err
	}

	return Sqrshl{
		rd:  rd.name(),
		rn:  rn.name(),
		rm:  rm.name(),
		arr: arr,
	}, nil
}

const sqrshlEnc uint32 = 237001728 // sqrshl vd, vn, vm (Q=0, size=0 form)

func (i Sqrshl) Encode(w io.Writer) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sqrshl: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("sqrshl: %w", err)
	}

	return writeWord(w, sqrshlEnc|q<<30|size<<22|rd|rn<<5|rm<<16)
}

func (i Sqrshl) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sqrshl.%s %s, %s, %s", i.arr, i.rd, i.rn, i.rm)
}
