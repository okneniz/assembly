package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AndImm — and rd, rn, #bitmask.
type AndImm struct {
	base
	logImm
}

const andImmX uint32 = 0x92000000

const andImmW uint32 = 0x12000000

// AndImm — and rd, rn, #bitmask. Register 31 reads as zr
// (SP/WSP are not allowed — use XZR/WZR); the mask must be encodable
// as a logical immediate (see encodeBitMasks).
func (Builder) AndImm(rd, rn Reg, imm uint64) (Instr, error) {
	if err := requireClass(rd, "AndImm", "rd", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AndImm", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("AndImm", rd, rn); err != nil {
		return nil, err
	}

	n, immr, imms, ok := encodeBitMasks(rd.Is64(), imm)
	if !ok {
		return nil, fmt.Errorf("arm64.NewAndImm: operand imm: %#x not encodable as bitmask", imm)
	}

	return AndImm{logImm: newLogImm(rd.name(), rn.name(), immr, imms, n == 1, rd.Is64())}, nil
}

func decodeAndImm(w uint32, addr uint64) Instr {
	return AndImm{
		newBase(addr, w),
		decodeLogImm(w),
	}
}

func (i AndImm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("and %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i AndImm) Encode(w io.Writer, pc uint64) (int64, error) {
	match := andImmX
	if !i.is64 {
		match = andImmW
	}

	if i.n {
		match |= 1 << 22
	}

	rd, rn, err := i.bits()
	if err != nil {
		return 0, fmt.Errorf("and: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|i.imms<<10|i.immr<<16)
}

func (i AndImm) MarshalJSON() ([]byte, error) {
	return i.marshal("and", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "immr": i.immr, "imms": i.imms})
}
