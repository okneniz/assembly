package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stur — stur ... (see lsBase for addressing kinds).
type Stur struct {
	base
	lsBase
}

// newSturBase - the Stur constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newSturBase(b base, e lsBase) Stur {
	return Stur{
		base:   b,
		lsBase: e,
	}
}

// newStur - the Stur constructor: validates the operands,
// delegates the assembly to newSturBase.
func newStur(b base, rt, rn Reg, off Off) (Stur, error) {
	if err := lsOperand(rt, rn, "Stur"); err != nil {
		return Stur{}, err
	}

	enc := sturWEnc
	if rt.Is64() {
		enc = sturXEnc
	}

	if err := requireUnscaledOff("Stur", off); err != nil {
		return Stur{}, err
	}

	return newSturBase(
		b,
		newLsBase(rt.name(), rn.name(), memUnscaled, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	sturXEnc uint32 = 0xF8000000 // stur xt, [xn, #±imm9]
	sturWEnc uint32 = 0xB8000000 // stur wt, [xn, #±imm9]
)

func (i Stur) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("stur %s, %s", i.rt, i.lsText(ctx))
}

func (i Stur) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "stur")
}
