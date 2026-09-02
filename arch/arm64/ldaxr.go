package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Ldaxr — ldaxr rt, [rn].
type Ldaxr struct {
	base
	atomic

	enc uint32
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	ldaxrXEnc uint32 = 0xC85FFC00 // ldaxr xt, [xn]
	ldaxrWEnc uint32 = 0x885FFC00 // ldaxr wt, [xn]
)

// Ldaxr — ldaxr rt, [rn]: rt — x/w register (register 31 reads as
// zr), rn — x register or SP (register 31 in the base reads as sp).
func (Builder) Ldaxr(rt, rn Reg) (Instr, error) {
	if err := lsOperand(rt, rn, "Ldaxr"); err != nil {
		return nil, err
	}

	enc := ldaxrWEnc
	if rt.Is64() {
		enc = ldaxrXEnc
	}

	return Ldaxr{
		atomic: newAtomic(rt.name(), rn.name()),
		enc:    enc,
	}, nil
}

func decodeLdaxrOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Ldaxr{
			base:   newBase(addr, w),
			atomic: newAtomic(armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:    enc,
		}
	}
}

func (i Ldaxr) ObjDump(_ disasm.ViewCtx) string {
	return "ldaxr " + i.atText()
}

func (i Ldaxr) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.atWrite(w, i.enc, "ldaxr")
}

func (i Ldaxr) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"ldaxr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rt": i.rt, "Rn": i.rn},
	)
}
