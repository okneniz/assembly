package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

const (
	movnX uint32 = 0x92800000
	movnW uint32 = 0x12800000
)

// Movn - movn rd, #(~(imm16 << hw*16)); displayed as mov.
type Movn struct {
	base

	rd        string
	imm16, hw uint32
}

func decodeMovn(w uint32, addr uint64) Instr {
	return Movn{
		base:  newBase(addr, w),
		rd:    armRegName(w&0x1f, w>>31&1 == 1),
		imm16: w >> 5 & 0xffff,
		hw:    w >> 21 & 0x3,
	}
}

func (i Movn) ObjDump(_ disasm.ViewCtx) string {
	is64 := i.rd[0] == 'x'
	if i.hw == 0 {
		return fmt.Sprintf("mov %s, #-0x%x", i.rd, uint64(i.imm16)+1)
	}

	shifted := uint64(i.imm16) << (i.hw * 16)
	val := ^shifted
	if !is64 {
		val &= 0xFFFFFFFF
	}

	if is64 && val&(uint64(1)<<63) != 0 {
		return fmt.Sprintf("mov %s, #-0x%x", i.rd, (^val)+1)
	}

	if !is64 && val&0x80000000 != 0 {
		return fmt.Sprintf("mov %s, #-0x%x", i.rd, (^val&0xFFFFFFFF)+1)
	}

	return fmt.Sprintf("mov %s, #0x%x", i.rd, val)
}

func (i Movn) Encode(w io.Writer, pc uint64) (int64, error) {
	match, err := sfMatch(i.rd, movnX, movnW)
	if err != nil {
		return 0, fmt.Errorf("movn: %w", err)
	}

	rd, err := armRegNum(i.rd)
	if err != nil {
		return 0, fmt.Errorf("movn: %w", err)
	}

	if i.imm16 > 0xffff || i.hw > 3 {
		return 0, errors.New("movn: imm/hw out of range")
	}

	return writeWord(w, match|rd|i.imm16<<5|i.hw<<21)
}

func (i Movn) MarshalJSON() ([]byte, error) {
	return i.marshal("movn", i.ObjDump(disasm.DefaultViewCtx()), "Data processing - immediate",
		map[string]any{"Rd": i.rd, "imm16": i.imm16, "hw": i.hw})
}
