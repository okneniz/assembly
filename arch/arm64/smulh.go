package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Smulh — smulh rd, rn, rm.
type Smulh struct {
	base

	rd, rn, rm string
}

const SmulhX uint32 = 0x9B407C00

func decodeSmulh(w uint32, addr uint64) Instr {
	return Smulh{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i Smulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("smulh %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Smulh) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, SmulhX, 0)
	if err != nil {
		return 0, fmt.Errorf("smulh: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("smulh: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Smulh) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"smulh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
