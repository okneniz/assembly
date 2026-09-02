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

// AndsImm — ands rd, rn, #bitmask (tst when Rd = zr). Register 31
// reads as zr (SP/WSP are not allowed — use XZR/WZR); the mask must be
// encodable as a logical immediate (see encodeBitMasks).
func (Builder) AndsImm(rd, rn Reg, imm uint64) (Instr, error) {
	if err := requireClass(
		rd,
		"AndsImm",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the tst form)",
		classX,
		classW,
		classXZR,
		classWZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(rn, "AndsImm", "rn", "register 31 reads as zr — use XZR/WZR",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if err := requireWidth("AndsImm", rd, rn); err != nil {
		return nil, err
	}

	n, immr, imms, ok := encodeBitMasks(rd.Is64(), imm)
	if !ok {
		return nil, fmt.Errorf("arm64.NewAndsImm: operand imm: %#x not encodable as bitmask", imm)
	}

	return AndsImm{logImm: newLogImm(rd.name(), rn.name(), immr, imms, n == 1, rd.Is64())}, nil
}

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
