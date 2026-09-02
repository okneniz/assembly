package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Cbnz — cbnz rt, target.
type Cbnz struct {
	base

	rt     string
	target imm
}

// Cbnz — cbnz rt, target: target — the absolute address of the
// branch destination (the ±1MB imm19 range is checked at encode time,
// from pc). rt — x/w register (register 31 reads as zr — use XZR/WZR).
func (Builder) Cbnz(rt Reg, target int64) (Instr, error) {
	if err := requireClass(rt, "Cbnz", "rt", "x/w register (register 31 reads as zr — use XZR/WZR)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	return Cbnz{rt: rt.name(), target: immNum(target)}, nil
}

func decodeCbnz(w uint32, addr uint64) Instr {
	return Cbnz{
		base:   newBase(addr, w),
		rt:     armRegName(w&0x1f, w>>31&1 == 1),
		target: immNum(int64(addr) + signExtendN(w>>5&0x7ffff, 19)*4),
	}
}

func (i Cbnz) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("cbnz %s, %s", i.rt, i.target.textHex())
}

func (i Cbnz) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	bits, err := brBits(target, int64(pc), 19)
	if err != nil {
		return 0, fmt.Errorf("cbnz: %w", err)
	}

	num, err := armRegNum(i.rt)
	if err != nil {
		return 0, fmt.Errorf("cbnz: %w", err)
	}

	match := uint32(0x35000000)
	if i.rt[0] == 'x' {
		match = 0xB5000000
	}

	return writeWord(w, match|bits<<5|num)
}

func (i Cbnz) MarshalJSON() ([]byte, error) {
	return i.marshal("cbnz", i.ObjDump(disasm.DefaultViewCtx()), "Branch",
		map[string]any{"Rt": i.rt, "imm19": i.target.val - int64(i.addr)})
}
