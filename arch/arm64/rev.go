package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rev — rev rd, rn.
type Rev struct {
	base

	rd, rn string
}

const RevX uint32 = 0xDAC00C00

func decodeRev(w uint32, addr uint64) Instr {
	return Rev{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
	}
}

func (i Rev) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev %s, %s", i.rd, i.rn)
}

func (i Rev) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, RevX, 0)
	if err != nil {
		return 0, fmt.Errorf("rev: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rev: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (i Rev) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"rev",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
