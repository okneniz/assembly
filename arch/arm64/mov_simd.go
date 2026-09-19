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

	// the mov alias of ORR: Rn == Rm (llvm encodes it so)
	return writeWord(w, i.enc|rd|rm<<5|rm<<16)
}

// SkipVerify — there is no decoding schema.
func (i MovSimd) SkipVerify() {}

// MovSimdBuilder entry: mov.8b/16b vd, vm (ORR-vector, Rn=31).
func (Builder) MovSimd(rd, rm VReg, arr string) (Instr, error) {
	err := requireArr("MovSimd", arr, "8b", "16b")
	if err != nil {
		return nil, err
	}

	enc := uint32(0x4EA01C00) // mov.16b (Q=1)
	if arr == "8b" {
		enc &^= 1 << 30
	}

	return MovSimd{
		rd:  rd.name(),
		rm:  rm.name(),
		arr: arr,
		enc: enc,
	}, nil
}
