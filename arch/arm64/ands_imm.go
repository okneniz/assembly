package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AndsImm — ands rd, rn, #bitmask; pseudo tst (Rd = zr).
type AndsImm struct {
	base
	logImm
}

const (
	andsImmX uint32 = 0xF2000000
	andsImmW uint32 = 0x72000000
)

func decodeAndsImm(w uint32, addr uint64) Instr {
	return AndsImm{
		newBase(addr, w),
		decodeLogImm(w),
	}
}

func (i AndsImm) ObjDump(_ disasm.ViewCtx) string {
	zr := "xzr"
	if !i.is64 {
		zr = "wzr"
	}

	if i.rd == zr {
		return fmt.Sprintf("tst %s, #0x%x", i.rn, i.mask())
	}

	return fmt.Sprintf("ands %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i AndsImm) Encode(w io.Writer, pc uint64) (int64, error) {
	match := andsImmX
	if !i.is64 {
		match = andsImmW
	}

	if i.n {
		match |= 1 << 22
	}

	rd, rn, err := i.bits()
	if err != nil {
		return 0, fmt.Errorf("ands: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|i.imms<<10|i.immr<<16)
}

func (i AndsImm) MarshalJSON() ([]byte, error) {
	return i.marshal("ands", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}
