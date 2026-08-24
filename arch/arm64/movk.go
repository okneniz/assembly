package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

const (
	movkX uint32 = 0xF2800000
	movkW uint32 = 0x72800000
)

// Movk — movk rd, #imm16[, lsl #hw*16].
type Movk struct {
	base

	rd        string
	imm16, hw uint32
}

// NewMovk - movk rd, #imm16, lsl #hw*16. The 32-bit form allows only
// Hw0/Hw1 (shift up to #16).
func NewMovk(rd Reg, imm Imm16, hw Hw) (Instr, error) {
	if err := requireClass(
		rd,
		"Movk",
		"rd",
		"x/w register, sp not allowed (register 31 reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	); err != nil {
		return nil, err
	}

	if err := requireHwW(rd, "Movk", hw); err != nil {
		return nil, err
	}

	return Movk{
		rd:    rd.name(),
		imm16: imm.v,
		hw:    uint32(hw),
	}, nil
}

func (i Movk) ObjDump(_ disasm.ViewCtx) string {
	if i.hw == 0 {
		return fmt.Sprintf("movk %s, #0x%x", i.rd, i.imm16)
	}

	return fmt.Sprintf("movk %s, #0x%x, lsl #%d", i.rd, i.imm16, i.hw*16)
}

func (i Movk) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, movkX, movkW)
	if err != nil {
		return 0, fmt.Errorf("movk: %w", err)
	}

	rd, err := armRegNum(i.rd)
	if err != nil {
		return 0, fmt.Errorf("movk: %w", err)
	}

	if i.imm16 > 0xffff || i.hw > 3 {
		return 0, errors.New("movk: imm/hw out of range")
	}

	return writeWord(w, match|rd|i.imm16<<5|i.hw<<21)
}

func (i Movk) MarshalJSON() ([]byte, error) {
	return i.marshal("movk", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "imm16": i.imm16, "hw": i.hw})
}

func decodeMovk(w uint32, addr uint64) Instr {
	return Movk{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		imm16: w >> 5 & 0xffff,
		hw:    w >> 21 & 0x3,
	}
}
