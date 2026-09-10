package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Movz — movz rd, #(imm16 << hw*16); displayed as mov.
type Movz struct {
	base

	rd        string
	imm16, hw uint32
}

// newMovz - the Movz constructor: validates the operands and
// assembles the struct (the Builder method delegates here; the
// decoder calls it with values read from the word).
func newMovz(b base, rd Reg, imm Imm16, hw Hw) (Movz, error) {
	err := requireClass(
		rd,
		"Movz",
		"rd",
		"x/w register, sp not allowed (register 31 reads as zr)",
		classX,
		classW,
		classXZR,
		classWZR,
	)

	if err != nil {
		return Movz{}, err
	}

	if err = requireHwW(rd, "Movz", hw); err != nil {
		return Movz{}, err
	}

	return Movz{
		base:  b,
		rd:    rd.name(),
		imm16: imm.v,
		hw:    uint32(hw),
	}, nil
}

const (
	movzX uint32 = 0xD2800000
	movzW uint32 = 0x52800000
)

func (i Movz) ObjDump(_ disasm.ViewCtx) string {
	is64 := i.rd[0] == 'x'
	val := uint64(i.imm16) << (i.hw * 16)
	if is64 && val&(uint64(1)<<63) != 0 {
		return fmt.Sprintf("mov %s, #-0x%x", i.rd, (^val)+1)
	}

	return fmt.Sprintf("mov %s, #0x%x", i.rd, val)
}

func (i Movz) Encode(w io.Writer) (int64, error) {
	match, err := sfMatch(i.rd, movzX, movzW)
	if err != nil {
		return 0, fmt.Errorf("movz: %w", err)
	}

	rd, err := armRegNum(i.rd)
	if err != nil {
		return 0, fmt.Errorf("movz: %w", err)
	}

	if i.imm16 > 0xffff || i.hw > 3 {
		return 0, errors.New("movz: imm/hw out of range")
	}

	return writeWord(w, match|rd|i.imm16<<5|i.hw<<21)
}

func (Builder) Movz(rd Reg, imm Imm16, hw Hw) (Instr, error) {
	return newMovz(base{}, rd, imm, hw)
}

func decodeMovz(w uint32) (Instr, error) {
	in, err := newMovz(
		newBase(w),
		gprOf(w&0x1f, w>>31&1 == 1),
		imm16Of(w>>5&0xffff),
		hwOf(w>>21&0x3),
	)
	if err != nil {
		return nil, err
	}

	return in, nil
}
