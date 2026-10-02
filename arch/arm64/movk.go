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
	rd        string
	imm16, hw uint32
}

// newMovk - the Movk constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newMovk(rd Reg, imm Imm16, hw Hw) (Movk, error) {
	err := requireClass(
		rd,
		"Movk",
		"rd",
		"x/w register, sp not allowed (register 31 reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Movk{}, err
	}

	if err = requireHwW(rd, "Movk", hw); err != nil {
		return Movk{}, err
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
