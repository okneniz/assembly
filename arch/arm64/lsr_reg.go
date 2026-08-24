package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// LsrReg — lsr rd, rn, rm.
type LsrReg struct {
	base

	rd, rn, rm string
}

const LsrRegX uint32 = 0x9A002400

func decodeLsrReg(w uint32, addr uint64) Instr {
	return LsrReg{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i LsrReg) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("lsr %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i LsrReg) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, LsrRegX, 0x1A002400)
	if err != nil {
		return 0, fmt.Errorf("lsr: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("lsr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i LsrReg) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"lsr",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
