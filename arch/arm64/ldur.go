package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldur — ldur ... (see lsBase for the addressing kinds).
type Ldur struct {
	base
	lsBase
}

// newLdurBase - the Ldur constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdurBase(b base, e lsBase) Ldur {
	return Ldur{
		base:   b,
		lsBase: e,
	}
}

// newLdur - the Ldur constructor: validates the operands,
// delegates the assembly to newLdurBase.
func newLdur(b base, rt, rn Reg, off Off) (Ldur, error) {
	if err := lsOperand(rt, rn, "Ldur"); err != nil {
		return Ldur{}, err
	}

	enc := ldurWEnc
	if rt.Is64() {
		enc = ldurXEnc
	}

	if err := requireUnscaledOff("Ldur", off); err != nil {
		return Ldur{}, err
	}

	return newLdurBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldurXEnc uint32 = 0xF8400000 // ldur xt, [xn, #±imm9]
	ldurWEnc uint32 = 0xB8400000 // ldur wt, [xn, #±imm9]
)

func (i Ldur) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldur %s, %s", i.rt, i.lsText(ctx))
}

func (i Ldur) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldur")
}
