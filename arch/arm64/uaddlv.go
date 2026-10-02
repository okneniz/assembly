package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Uaddlv — uaddlv.Arr hN/sN/dN, vn (scalar dest by size).
type Uaddlv struct {
	rd, rn  string // rd - the scalar print name (hN/sN/dN)
	q, size uint32
}

// newUaddlv - the Uaddlv constructor: the struct is assembled only
// here (the Builder method and the decoder call it); the destination
// print name is the scalar view of the register number by size.
func newUaddlv(q, size uint32, rd, rn VReg) (Uaddlv, error) {
	scalar := fmt.Sprintf("d%d", rd.Num())
	if size == 0 {
		scalar = fmt.Sprintf("h%d", rd.Num())
	} else if size == 1 {
		scalar = fmt.Sprintf("s%d", rd.Num())
	}

	return Uaddlv{
		rd:   scalar,
		rn:   rn.name(),
		q:    q,
		size: size,
	}, nil
}

const uaddlvEnc uint32 = 0x2E303800

func (i Uaddlv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("uaddlv.%s %s, %s", decodeArrangement(i.q, i.size), i.rd, i.rn)
}

func (i Uaddlv) Encode(w io.Writer) (int64, error) {
	rd, rn, err := regNums2(fmt.Sprintf("v%d", regIndex(i.rd)), i.rn)
	if err != nil {
		return 0, fmt.Errorf("uaddlv: %w", err)
	}

	return writeWord(w, uaddlvEnc|rd|rn<<5)
}
