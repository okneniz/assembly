package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// OrrImm — orr rd, rn, #bitmask.
type OrrImm struct {
	base
	logImm
}

// newOrrImm - the OrrImm constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newOrrImm(b base, rd Reg, rn Reg, imm uint64) (OrrImm, error) {
	err := requireClass(
		rd,
		"OrrImm",
		"rd",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return OrrImm{}, err
	}

	err = requireClass(
		rn,
		"OrrImm",
		"rn",
		"register 31 reads as zr — use XZR/WZR",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return OrrImm{}, err
	}

	err = requireWidth(
		"OrrImm",
		rd,
		rn,
	)

	if err != nil {
		return OrrImm{}, err
	}

	n, immr, imms, ok := encodeBitMasks(rd.Is64(), imm)
	if !ok {
		return OrrImm{}, fmt.Errorf(
			"arm64.NewOrrImm: operand imm: %#x not encodable as bitmask",
			imm,
		)
	}

	return OrrImm{
		base:   b,
		logImm: newLogImm(rd.name(), rn.name(), immr, imms, n == 1, rd.Is64()),
	}, nil
}

const (
	orrImmX uint32 = 0xB2000000
	orrImmW uint32 = 0x32000000
)

// movzRep/movnRep - whether the pattern is representable by a single
// MOVZ/MOVN (set bits / in the MOVN case zero bits fit into a single
// 16-bit aligned window).
func movzRep(v uint64, bits int) bool {
	for s := uint(0); s < uint(bits); s += 16 {
		if v>>s&0xffff != 0 && v&^(0xffff<<s) == 0 {
			return true
		}
	}

	return false
}

func movnRep(v uint64, bits int) bool {
	all := ^uint64(0)
	if bits == 32 {
		all = 0xffffffff
	}

	return movzRep(^v&all, bits)
}

// ObjDump — orr rd, rn, #bitmask. With Rn = ZR LLVM objdump prints the mov
// alias, but ONLY if the immediate is not representable by a single
// MOVZ/MOVN (otherwise the text is ambiguous on reassembly:
// mov x, #0x20 canonically = the MOVZ encoding).
func (i OrrImm) ObjDump(_ disasm.ViewCtx) string {
	bits := 32
	if i.is64 {
		bits = 64
	}

	if (i.rn == "xzr" || i.rn == "wzr") &&
		!movzRep(i.mask(), bits) && !movnRep(i.mask(), bits) {
		return fmt.Sprintf("mov %s, %s", i.rd, i.immText())
	}

	return fmt.Sprintf("orr %s, %s, #0x%x", i.rd, i.rn, i.mask())
}

func (i OrrImm) Encode(w io.Writer) (int64, error) {
	match := orrImmX
	if !i.is64 {
		match = orrImmW
	}

	if i.n {
		match |= 1 << 22
	}

	rd, rn, err := i.bits()
	if err != nil {
		return 0, fmt.Errorf("orr: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|i.imms<<10|i.immr<<16)
}

// immText - the immediate in objdump style: 64-bit patterns with the top
// bit set are printed signed (#-0x80000000 instead of #0xffffffff80000000).
func (i OrrImm) immText() string {
	m := i.mask()
	if i.is64 && m>>63 == 1 {
		return fmt.Sprintf("#-0x%x", ^m+1)
	}

	return fmt.Sprintf("#0x%x", m)
}

func (Builder) OrrImm(rd, rn Reg, imm uint64) (Instr, error) {
	return newOrrImm(base{}, rd, rn, imm)
}

func decodeOrrImm(w uint32) (Instr, error) {
	in, err := newOrrImm(
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
