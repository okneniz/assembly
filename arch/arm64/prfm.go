package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Prfm — prfm pldl1keep, [rn]. The pldl1keep operand is a fixed keyword
// (not a register/immediate), so self-verify against renderInstr (which
// resolves Sym into numbers) is impossible - skipVerify.
type Prfm struct {
	base

	rn string
}

// Prfm — prfm pldl1keep, [rn]: the prefetch op is fixed (this type
// encodes the one form); rn — x register or SP (register 31 in the base
// reads as sp).
func (Builder) Prfm(rn Reg) (Instr, error) {
	if err := requireClass(
		rn,
		"Prfm",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	); err != nil {
		return nil, err
	}

	return Prfm{rn: rn.name()}, nil
}

func decodePrfm(w uint32, addr uint64) Instr {
	return Prfm{
		base: newBase(addr, w),
		rn:   regNameXSP(w >> 5 & 0x1f),
	}
}

func (i Prfm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("prfm pldl1keep, [%s]", i.rn)
}

func (i Prfm) Encode(w io.Writer, pc uint64) (int64, error) {
	n, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("prfm: %w", err)
	}

	return writeWord(w, 0xF9800000|n<<5)
}

func (i Prfm) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"prfm",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Load/Store",
		map[string]any{"Rn": i.rn},
	)
}

// SkipVerify — pldl1keep is a keyword, not an address.
func (i Prfm) SkipVerify() {}
