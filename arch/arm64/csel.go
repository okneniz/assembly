package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Csel — csel rd, rn, rm, cond. Csinc/Csinv/Csneg — the same encodings
// with the inverse aliases cset/csetm/cinc/cinv/cneg.
type Csel struct {
	base

	rd, rn, rm string
	cond       string
}

func NewCsel(rd string, rn string, rm string, cond string) Csel {
	return Csel{
		rd:   rd,
		rn:   rn,
		rm:   rm,
		cond: cond,
	}
}

const (
	cselX uint32 = 0x9A800000
	cselW uint32 = 0x1A800000
)

func decodeCsel(w uint32, addr uint64) Instr {
	return Csel{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
		cond: condName(w >> 12 & 0xf),
	}
}

func (i Csel) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("csel %s, %s, %s, %s", i.rd, i.rn, i.rm, i.cond)
}

func (i Csel) Encode(w io.Writer, pc uint64) (int64, error) {
	return cselWrite(w, i, cselX, cselW, "csel")
}

func (i Csel) MarshalJSON() ([]byte, error) {
	return i.marshal("csel", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "cond": i.cond})
}

// cselWrite — the common encoding skeleton of the csel family.
func cselWrite(w io.Writer, i Csel, matchX, matchW uint32, name string) (int64, error) {
	match, err := sfMatch(i.rd, matchX, matchW)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	c, err := condNum(i.cond)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}

	return writeWord(w, match|rd|rn<<5|c<<12|rm<<16)
}
