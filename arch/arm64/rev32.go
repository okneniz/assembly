package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Rev32 — rev32 rd, rn.
type Rev32 struct {
	base

	rd, rn string
}

const Rev32X uint32 = 0xDAC00800

// Rev32 — rev32 rd, rn. Only the 64-bit form: the architecture has
// no 32-bit rev32 (the sf=0 slot of the encoding is the 32-bit rev).
func (Builder) Rev32(rd, rn Reg) (Instr, error) {
	if err := requireClass(
		rd,
		"Rev32",
		"rd",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	if err := requireClass(
		rn,
		"Rev32",
		"rn",
		"only x registers (X/XZR)",
		classX,
		classXZR,
	); err != nil {
		return nil, err
	}

	return Rev32{
		rd: rd.name(),
		rn: rn.name(),
	}, nil
}

func decodeRev32(w uint32, addr uint64) Instr {
	return Rev32{
		base: newBase(addr, w),
		rd:   armRegName(w&0x1f, w>>31&1 == 1),
		rn:   armRegName(w>>5&0x1f, w>>31&1 == 1),
	}
}

func (i Rev32) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("rev32 %s, %s", i.rd, i.rn)
}

func (i Rev32) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, Rev32X, 0x5AC00800)
	if err != nil {
		return 0, fmt.Errorf("rev32: %w", err)
	}

	rd, rn, err := regNums2(i.rd, i.rn)
	if err != nil {
		return 0, fmt.Errorf("rev32: %w", err)
	}

	return writeWord(w, match|rd|rn<<5)
}

func (i Rev32) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"rev32",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing",
		map[string]any{"Rd": i.rd, "Rn": i.rn},
	)
}
