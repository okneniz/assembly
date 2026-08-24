package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Sdiv — sdiv rd, rn, rm.
type Sdiv struct {
	base

	rd, rn, rm string
}

const SdivX uint32 = 0x9AC00C00

func decodeSdiv(w uint32, addr uint64) Instr {
	return Sdiv{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i Sdiv) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sdiv %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Sdiv) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, SdivX, 0x1AC00C00)
	if err != nil {
		return 0, fmt.Errorf("sdiv: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("sdiv: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Sdiv) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"sdiv",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
