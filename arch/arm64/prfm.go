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

// newPrfm - the Prfm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newPrfm(b base, rn Reg) (Prfm, error) {
	err := requireClass(
		rn,
		"Prfm",
		"rn",
		"x register or SP (register 31 in the base reads as sp)",
		classX,
		classSP,
	)

	if err != nil {
		return Prfm{}, err
	}

	return Prfm{
		base: b,
		rn:   rn.name(),
	}, nil
}

func (i Prfm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("prfm pldl1keep, [%s]", i.rn)
}

func (i Prfm) Encode(w io.Writer) (int64, error) {
	n, err := armRegNum(i.rn)
	if err != nil {
		return 0, fmt.Errorf("prfm: %w", err)
	}

	return writeWord(w, 0xF9800000|n<<5)
}

// SkipVerify — pldl1keep is a keyword, not an address.
func (i Prfm) SkipVerify() {}

func (Builder) Prfm(rn Reg) (Instr, error) {
	return newPrfm(base{}, rn)
}

func decodePrfm(w uint32) (Instr, error) {
	in, err := newPrfm(newBase(w), xspOf(w>>5&0x1f))
	if err != nil {
		return nil, err
	}

	return in, nil
}
