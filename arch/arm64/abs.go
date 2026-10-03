package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Abs — abs.Arr vd, vn (absolute value per lane).
type Abs struct {
	rd, rn string
	arr    string // 8b/16b
}

// newAbs - the Abs constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the decoder
// calls it with values read from the word).
func newAbs(rd, rn VReg, arr string) (Abs, error) {
	err := requireArr("Abs", arr, "8b", "16b", "4h", "8h", "2s", "4s", "2d")
	if err != nil {
		return Abs{}, err
	}

	return Abs{
		rd:  rd.name(),
		rn:  rn.name(),
		arr: arr,
	}, nil
}

const absEnc uint32 = 0x0E20B800 // abs vd, vn (Q=0 form)

func (i Abs) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("abs: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("abs: %w", err)
	}

	return writeWord(w, absEnc|q<<30|size<<22|rd|rn<<5)
}

func (i Abs) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("abs.%s %s, %s", i.arr, i.rd, i.rn)
}
