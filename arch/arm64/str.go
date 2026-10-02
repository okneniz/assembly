package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Str — str ... (see lsBase for addressing kinds).
type Str struct {
	lsBase
}

// newStrBase - the Str constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newStrBase(e lsBase) Str {
	return Str{
		lsBase: e,
	}
}

// newStr - the Str constructor: validates the operands,
// delegates the assembly to newStrBase.
func newStr(rt, rn Reg, off Off) (Str, error) {
	if err := lsOperand(rt, rn, "Str"); err != nil {
		return Str{}, err
	}

	enc, scale := strXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = strWEnc, 2
	}

	if err := requireOff("Str", off, scale); err != nil {
		return Str{}, err
	}

	return newStrBase(
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	), nil
}

// Encodings of the unsigned-offset form: the access size is set by rt, the
// offset scale = log2 of the size.
const (
	strXEnc uint32 = 0xF9000000 // str xt, [xn, #imm12<<3]
	strWEnc uint32 = 0xB9000000 // str wt, [xn, #imm12<<2]
)

func (i Str) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("str %s, %s", i.rt, i.lsText(ctx))
}

func (i Str) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "str")
}

// The FP unsigned-offset encodings: the access size is set by the rt
// kind, the offset scale = log2 of the size.
const (
	strFDEnc uint32 = 0xFD000000 // str dt, [xn, #imm12<<3]
	strFSEnc uint32 = 0xBD000000 // str st, [xn, #imm12<<2]
)
