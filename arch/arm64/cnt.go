package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cnt — cnt.8b/16b vd, vn (population count per byte lane).
type Cnt struct {
	base

	rd, rn string
	arr    string // 8b/16b
}

// newCnt - the Cnt constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newCnt(b base, rd, rn VReg, arr string) (Cnt, error) {
	err := requireArr("Cnt", arr, "8b", "16b")
	if err != nil {
		return Cnt{}, err
	}

	return Cnt{
		base: b,
		rd:   rd.name(),
		rn:   rn.name(),
		arr:  arr,
	}, nil
}

const cntEnc uint32 = 0x0E205800 // cnt vd, vn (Q=0 form)

func (i Cnt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cnt.%s %s, %s", i.arr, i.rd, i.rn)
}

func (i Cnt) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("cnt: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("cnt: %w", err)
	}

	return writeWord(w, cntEnc|q<<30|size<<22|rd|rn<<5)
}
