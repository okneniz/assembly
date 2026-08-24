package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Madd - madd rd, rn, rm, ra; pseudo: mul (ra = xzr). Msub/msub/mneg are
// the same encodings with o0=1.
type Madd struct {
	base

	rd, rn, rm, ra string
}

func NewMadd(rd string, rn string, rm string, ra string) Madd {
	return Madd{
		rd: rd,
		rn: rn,
		rm: rm,
		ra: ra,
	}
}

const (
	maddX uint32 = 0x9B000000
	maddW uint32 = 0x1B000000
)

func decodeMadd(w uint32, addr uint64) Instr {
	return Madd{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		ra:   armRegName(w>>10&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i Madd) ObjDump(_ disasm.ViewCtx) string {
	zr := "xzr"
	if i.rd[0] == 'w' {
		zr = "wzr"
	}

	if i.ra == zr {
		return fmt.Sprintf("mul %s, %s, %s", i.rd, i.rn, i.rm)
	}

	return fmt.Sprintf("madd %s, %s, %s, %s", i.rd, i.rn, i.rm, i.ra)
}

func (i Madd) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, maddX, maddW)
	if err != nil {
		return 0, fmt.Errorf("madd: %w", err)
	}

	return msubWrite(w, match, i)
}

func (i Madd) MarshalJSON() ([]byte, error) {
	return i.marshal("madd", i.ObjDump(disasm.DefaultViewCtx()), "Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm, "Ra": i.ra})
}

// msubWrite - the shared word of the madd/msub family (msub = Madd with the opcode bit).
func msubWrite(w io.Writer, match uint32, i Madd) (int64, error) {
	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, err
	}

	ra, err := armRegNum(i.ra)
	if err != nil {
		return 0, err
	}

	return writeWord(w, match|rd|rn<<5|ra<<10|rm<<16)
}
