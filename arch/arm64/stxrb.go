package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Stxrb — stxrb rs, rt, [rn].
type Stxrb struct {
	base
	excl

	enc uint32
}

const stxrbEnc uint32 = 0x08000000 // stxrb ws, wt, [xn]

// Stxrb — stxrb rs, rt, [rn]: byte access, rs — the w status
// register (register 31 reads as wzr), rt — w register only (register 31
// reads as wzr), rn — x register or SP (register 31 in the base reads
// as sp).
func (Builder) Stxrb(rs, rt, rn Reg) (Instr, error) {
	if err := requireClass(rs, "Stxrb", "rs", "w status register (register 31 in rs reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rt, "Stxrb", "rt", "w register (register 31 in rt reads as wzr)",
		classW, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Stxrb",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	return Stxrb{
		excl: newExcl(rs.name(), rt.name(), rn.name()),
		enc:  stxrbEnc,
	}, nil
}

func decodeStxrbOf(enc uint32, x64 bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return Stxrb{
			base: newBase(addr, w),
			excl: newExcl(regNameW(w>>16&0x1f), armRegName(w&0x1f, x64), regNameXSP(w>>5&0x1f)),
			enc:  enc,
		}
	}
}

func (i Stxrb) ObjDump(_ disasm.ViewCtx) string {
	return "stxrb " + i.exText()
}

func (i Stxrb) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.exWrite(w, i.enc, "stxrb")
}

func (i Stxrb) MarshalJSON() ([]byte, error) {
	return i.marshal("stxrb", i.ObjDump(disasm.DefaultViewCtx()), "Load/Store",
		map[string]any{"Rs": i.rs, "Rt": i.rt, "Rn": i.rn})
}
