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

// newEorImm - the EorImm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newEorImm(b base, rd Reg, rn Reg, imm uint64) (EorImm, error) {
	err := requireClass(
		rd,
		"EorImm",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EorImm{}, err
	}

	err = requireClass(
		rn,
		"EorImm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return EorImm{}, err
	}

	err = requireWidth(
		"EorImm",
		rd,
		rn,
	)

	if err != nil {
		return EorImm{}, err
	}

	n, immr, imms, ok := encodeBitMasks(rd.Is64(), imm)
	if !ok {
		return EorImm{}, fmt.Errorf(
			"arm64.NewEorImm: operand imm: %#x not encodable as bitmask",
			imm,
		)
	}

	return EorImm{
		base:   b,
		logImm: newLogImm(rd.name(), rn.name(), immr, imms, n == 1, rd.Is64()),
	}, nil
}

const (
	eorImmX uint32 = 0xD2000000
	eorImmW uint32 = 0x52000000
)

func (i EorImm) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("eor %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i EorImm) Encode(w io.Writer) (int64, error) {
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

func (Builder) EorImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newEorImm(base{}, rd, rn, imm)
}

func decodeEorImm(w uint32) (Instr, error) {
	in, err := newEorImm(
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
