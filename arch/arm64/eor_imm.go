package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// EorImm — eor rd, rn, #bitmask.
type EorImm struct {
	base
	logImm
}

const (
	eorImmX uint32 = 0xD2000000
	eorImmW uint32 = 0x52000000
)

func decodeEorImm(w uint32, addr uint64) Instr {
	return EorImm{
		newBase(addr, w),
		decodeLogImm(w),
	}
}

func (i EorImm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("eor %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i EorImm) Encode(w io.Writer, pc uint64) (int64, error) {
	match := eorImmX
	if !i.is64 {
		match = eorImmW
	}

	if i.n {
		match |= 1 << 22
	}

	rd, rn, err := i.bits()
	if err != nil {
		return 0, fmt.Errorf("eor: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|i.imms<<10|i.immr<<16)
}

func (i EorImm) MarshalJSON() ([]byte, error) {
	return i.marshal("eor", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}
