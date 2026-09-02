package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlxr — stlxr rs, rt, [rn].
type Stlxr struct {
	base
	excl

	enc uint32
}

// Encodings of the 64/32-bit forms: the access size is set by rt.
const (
	stlxrXEnc uint32 = 0xC800FC00 // stlxr ws, xt, [xn]
	stlxrWEnc uint32 = 0x8800FC00 // stlxr ws, wt, [xn]
)

// Stlxr — stlxr rs, rt, [rn]: rs — the w status register
// (register 31 reads as wzr), rt — x/w register (register 31 reads as
// zr), rn — x register or SP (register 31 in the base reads as sp).
func (Builder) Stlxr(rs, rt, rn Reg) (Instr, error) {
	if err := requireClass(rs, "Stlxr", "rs", "w status register (register 31 in rs reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := lsOperand(rt, rn, "Stlxr"); err != nil {
		return nil, err
	}

	enc := stlxrWEnc
	if rt.Is64() {
		enc = stlxrXEnc
	}

	return Stlxr{
		excl: newExcl(rs.name(), rt.name(), rn.name()),
		enc:  enc,
	}, nil
}

func decodeStlxrOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stlxr{
			base: newBase(addr, w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}
	}
}

func (i Stlxr) ObjDump(_ disasm.ViewCtx) string {
	return "stlxr " + i.exText()
}

func (i Stlxr) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.exWrite(w, i.enc, "stlxr")
}

func (i Stlxr) MarshalJSON() ([]byte, error) {
	return i.marshal("stlxr", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rs": i.rs, "Rt": i.rt, "Rn": i.rn})
}
