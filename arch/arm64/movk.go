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

// newMovk - the Movk constructor: the struct is assembled only
// here (the Builder method and the decoder call it).
func newMovk(b base, rd string, imm16 uint32, hw uint32) Movk {
	return Movk{
		base:  b,
		rd:    rd,
		imm16: imm16,
		hw:    hw,
	}
}

func (i Movk) ObjDump(_ disasm.ViewCtx) string {
	if i.hw == 0 {
		return fmt.Sprintf("movk %s, #0x%x", i.rd, i.imm16)
	}

	return fmt.Sprintf("movk %s, #0x%x, lsl #%d", i.rd, i.imm16, i.hw*16)
}

func (i Movk) Encode(w io.Writer) (int64, error) {
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

// Movk — movk rd, #imm16, lsl #hw*16. The 32-bit form allows only
// Hw0/Hw1 (shift up to #16).
func (Builder) Movk(rd Reg, imm Imm16, hw Hw) (Instr, error) {
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

	return newMovk(base{}, rd.name(), imm.v, uint32(hw)), nil
}

func decodeMovk(w uint32) Instr {
	return newMovk(newBase(w), armRegName(w&0x1f, w>>31&1 == 1), w>>5&0xffff, w>>21&0x3)
}
