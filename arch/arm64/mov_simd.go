package arm64

// MovSimd — mov.16b/mov.8b vd, vm: ORR-vector with Rn == Rm (the mov
// alias). The decoder does not know this encoding (no schema) - skipVerify:
// the constructor's encoding is unambiguous.

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

type MovSimd struct {
	base

	rd, rm string
	arr    string
	enc    uint32
}

func (i MovSimd) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("mov.%s %s, %s", i.arr, i.rd, i.rm)
}

func (i MovSimd) Encode(w io.Writer) (int64, error) {
	rd, rm, err := regNums2(i.rd, i.rm)
	if err != nil {
		return 0, fmt.Errorf("mov: %w", err)
	}

	return writeWord(w, i.enc|rd|31<<5|rm<<16)
}

// SkipVerify — there is no decoding schema.
func (i MovSimd) SkipVerify() {}
