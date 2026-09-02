package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlr — stlr rt, [rn].
type Stlr struct {
	base
	atomic

	enc uint32
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	stlrXEnc uint32 = 0xC89FFC00 // stlr xt, [xn]
	stlrWEnc uint32 = 0x889FFC00 // stlr wt, [xn]
)

// Stlr — stlr rt, [rn]: rt — x/w register (register 31 reads as
// zr), rn — x register or SP (register 31 in the base reads as sp).
func (Builder) Stlr(rt, rn Reg) (Instr, error) {
	if err := lsOperand(rt, rn, "Stlr"); err != nil {
		return nil, err
	}

	enc := stlrWEnc
	if rt.Is64() {
		enc = stlrXEnc
	}

	return Stlr{
		atomic: newAtomic(rt.name(), rn.name()),
		enc:    enc,
	}, nil
}

func decodeStlrOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stlr{
			base:   newBase(addr, w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}
	}
}

func (i Stlr) ObjDump(_ disasm.ViewCtx) string {
	return "stlr " + i.atText()
}

func (i Stlr) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.atWrite(w, i.enc, "stlr")
}

func (i Stlr) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"stlr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
