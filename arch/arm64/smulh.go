package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Smulh — smulh rd, rn, rm.
type Smulh struct {
	base

	rd, rn, rm string
}

// newSmulh - the Smulh constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newSmulh(b base, rd string, rn string, rm string) Smulh {
	return Smulh{
		base: b,
		rd:   rd,
		rn:   rn,
		rm:   rm,
	}
}

const SmulhX uint32 = 0x9B407C00

func (i Smulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("smulh %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Smulh) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, SmulhX, 0)
	if err != nil {
		return 0, fmt.Errorf("smulh: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("smulh: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

// Smulh — smulh rd, rn, rm. Only the 64-bit form (the architecture
// has no 32-bit smulh); register 31 reads as zr (use XZR).
func (Builder) Smulh(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(
		rd,
		"Smulh",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Smulh",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rm,
		"Smulh",
		"rm",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return newSmulh(base{}, rd.name(), rn.name(), rm.name()), nil
}

func decodeSmulh(w uint32) Instr {
	return newSmulh(
		newBase(w),
		armRegName(w&0x1f, w>>31&1 == 1),
		armRegName(w>>5&0x1f, w>>31&1 == 1),
		armRegName(w>>16&0x1f, w>>31&1 == 1),
	)
}
