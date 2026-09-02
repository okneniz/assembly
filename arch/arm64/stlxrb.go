package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stlxrb — stlxrb rs, rt, [rn].
type Stlxrb struct {
	base
	excl

	enc uint32
}

const stlxrbEnc uint32 = 0x0800FC00 // stlxrb ws, wt, [xn]

// Stlxrb — stlxrb rs, rt, [rn]: byte access, rs — the w status
// register (register 31 reads as wzr), rt — w register only (register 31
// reads as wzr), rn — x register or SP (register 31 in the base reads
// as sp).
func (Builder) Stlxrb(rs, rt, rn Reg) (Instr, error) {
	if err := requireClass(rs, "Stlxrb", "rs", "w status register (register 31 in rs reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rt, "Stlxrb", "rt", "w register (register 31 in rt reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Stlxrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	return Stlxrb{
		excl: newExcl(rs.name(), rt.name(), rn.name()),
		enc:  stlxrbEnc,
	}, nil
}

func decodeStlxrbOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stlxrb{
			base: newBase(addr, w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}
	}
}

func (i Stlxrb) ObjDump(_ disasm.ViewCtx) string {
	return "stlxrb " + i.exText()
}

func (i Stlxrb) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.exWrite(w, i.enc, "stlxrb")
}

func (i Stlxrb) MarshalJSON() ([]byte, error) {
	return i.marshal("stlxrb", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rs": i.rs, "Rt": i.rt, "Rn": i.rn})
}
