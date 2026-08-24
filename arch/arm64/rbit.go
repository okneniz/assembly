package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rbit — rbit rd, rn.
type Rbit struct {
	base

	rd, rn string
}

const RbitX uint32 = 0xDAC00000

func decodeRbit(w uint32, addr uint64) Instr {
	return Rbit{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
	}
}

func (i Rbit) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rbit %s, %s", i.rd, i.rn)
}

func (i Rbit) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, RbitX, 0x5AC00000)
	if err != nil {
		return 0, fmt.Errorf("rbit: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rbit: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (i Rbit) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"rbit",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
