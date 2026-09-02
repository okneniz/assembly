package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Umulh — umulh rd, rn, rm.
type Umulh struct {
	base

	rd, rn, rm string
}

const UmulhX uint32 = 0x9BC07C00

// Umulh — umulh rd, rn, rm. Only the 64-bit form (the architecture
// has no 32-bit umulh); register 31 reads as zr (use XZR).
func (Builder) Umulh(rd, rn, rm Reg) (Instr, error) {
	if err := requireClass(
		rd,
		"Umulh",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Umulh",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rm,
		"Umulh",
		"rm",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return Umulh{
		rd: rd.name(),
		rn: rn.name(),
		rm: rm.name(),
	}, nil
}

func decodeUmulh(w uint32, addr uint64) Instr {
	return Umulh{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
		rm:   armRegName(w>>16&0x1f, w>>31&1 == 1),
	}
}

func (i Umulh) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("umulh %s, %s, %s", i.rd, i.rn, i.rm)
}

func (i Umulh) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, UmulhX, 0)
	if err != nil {
		return 0, fmt.Errorf("umulh: %w", err)
	}

	rd, rn, rm, err := regNums3(i.rd, i.rn, i.rm)
	if err != nil {
		return 0, fmt.Errorf("umulh: %w", err)
	}

	return writeWord(w, match|rd|rn<<5|rm<<16)
}

func (i Umulh) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"umulh",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn, "Rm": i.rm},
	)
}
