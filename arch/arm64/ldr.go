package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldr — ldr ... (see lsBase for the addressing kinds).
type Ldr struct {
	lsBase
}

// Encodings of the unsigned-offset form: the access size is set by rt, the
// offset scale = log2 of the size.
const (
	ldrXEnc uint32 = 0xF9400000 // ldr xt, [xn, #imm12<<3]
	ldrWEnc uint32 = 0xB9400000 // ldr wt, [xn, #imm12<<2]
)

func (i Ldr) Encode(w io.Writer) (int64, error) {
	return i.lsWrite(w, "ldr")
}

func (i Ldr) ObjDump(ctx disasm.ViewCtx) string {
	return fmt.Sprintf("ldr %s, %s", i.rt, i.lsText(ctx))
}

// ldrPoolWrap — pool ldr without self-verify: the decoder prints the slot's
// absolute address, not "=literal" — the text is not reproducible; the
// encoding is unambiguous (imm19 from the slot's known address).
type ldrPoolWrap struct {
	Ldr
}

func (ldrPoolWrap) SkipVerify() {}

// LdrPoolWrapOf — the pool-wrapped literal ldr (PoolUser); lit — the
// pc-relative byte offset of the pool slot.
func LdrPoolWrapOf(rt string, lit int64, enc uint32) Instr {
	return ldrPoolWrap{Ldr{lsBase: newLsBase(rt, "", memLiteral, 0, lit, enc, "", "", 0)}}
}

// newLdrBase - the Ldr constructor for a ready embedded base
// (the decoder and the string-operand layer): the struct is
// assembled only here.
func newLdrBase(e lsBase) Ldr {
	return Ldr{
		lsBase: e,
	}
}

// newLdr - the Ldr constructor: validates the operands,
// delegates the assembly to newLdrBase.
func newLdr(rt, rn Reg, off Off) (Ldr, error) {
	if err := lsOperand(rt, rn, "Ldr"); err != nil {
		return Ldr{}, err
	}

	enc, scale := ldrXEnc, uint32(3)
	if !rt.Is64() {
		enc, scale = ldrWEnc, 2
	}

	if err := requireOff("Ldr", off, scale); err != nil {
		return Ldr{}, err
	}

	return newLdrBase(
		newLsBase(rt.name(), rn.name(), memImm, int64(off), 0, enc, "", "", 0),
	), nil
}

// The FP unsigned-offset encodings: the access size is set by the rt
// kind, the offset scale = log2 of the size.
const (
	ldrFDEnc uint32 = 0xFD400000 // ldr dt, [xn, #imm12<<3]
	ldrFSEnc uint32 = 0xBD400000 // ldr st, [xn, #imm12<<2]
)
