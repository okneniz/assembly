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

// newAndsImm - the AndsImm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAndsImm(b base, rd Reg, rn Reg, imm uint64) (AndsImm, error) {
	err := requireClass(
		rd,
		"AndsImm",
		"rd",
		"register 31 reads as zr — use XZR/WZR (the tst form)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndsImm{}, err
	}

	err = requireClass(
		rn,
		"AndsImm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndsImm{}, err
	}

	err = requireWidth(
		"AndsImm",
		rd,
		rn,
	)

	if err != nil {
		return AndsImm{}, err
	}

	n, immr, imms, ok := encodeBitMasks(rd.Is64(), imm)
	if !ok {
		return AndsImm{}, fmt.Errorf(
			"arm64.NewAndsImm: operand imm: %#x not encodable as bitmask",
			imm,
		)
	}

	return AndsImm{
		base:   b,
		logImm: newLogImm(rd.name(), rn.name(), immr, imms, n == 1, rd.Is64()),
	}, nil
}

const (
	andsImmX uint32 = 0xF2000000
	andsImmW uint32 = 0x72000000
)

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

func (i AndsImm) Encode(w io.Writer) (int64, error) {
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

func (Builder) AndsImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newAndsImm(base{}, rd, rn, imm)
}

func decodeAndsImm(w uint32) (Instr, error) {
	in, err := newAndsImm(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		gprOf(w>>5&0x1f, w>>31&1 == 1),
		decodeBitMasks(w>>22&1 == 1, w>>16&0x3f, w>>10&0x3f, w>>31&1 == 1),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
