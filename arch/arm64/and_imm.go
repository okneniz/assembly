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

// newAndImm - the AndImm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newAndImm(b base, rd Reg, rn Reg, imm uint64) (AndImm, error) {
	err := requireClass(
		rd,
		"AndImm",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndImm{}, err
	}

	err = requireClass(
		rn,
		"AndImm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return AndImm{}, err
	}

	err = requireWidth(
		"AndImm",
		rd,
		rn,
	)

	if err != nil {
		return AndImm{}, err
	}

	n, immr, imms, ok := encodeBitMasks(rd.Is64(), imm)
	if !ok {
		return AndImm{}, fmt.Errorf(
			"arm64.NewAndImm: operand imm: %#x not encodable as bitmask",
			imm,
		)
	}

	return AndImm{
		base:   b,
		logImm: newLogImm(rd.name(), rn.name(), immr, imms, n == 1, rd.Is64()),
	}, nil
}

const andImmX uint32 = 0x92000000

const andImmW uint32 = 0x12000000

func (i AndImm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("and %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i AndImm) Encode(w io.Writer) (int64, error) {
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

func (Builder) AndImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newAndImm(base{}, rd, rn, imm)
}

func decodeAndImm(w uint32) (Instr, error) {
	in, err := newAndImm(
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
