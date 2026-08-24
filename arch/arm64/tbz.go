package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tbz - tbz rt, #bit, target (b5 selects the x/w width of Rt).
type Tbz struct {
	base

	rt     string
	bit    uint32
	target imm
	isTbnz bool
}

func decodeTbzOf(isTbnz bool) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		x64 := w>>31&1 == 1
		return Tbz{
			base:   newBase(addr, w),
			rt:     armRegName(w&0x1f, x64),
			bit:    w>>19&0x1f | w>>26&0x20,
			target: immNum(int64(addr) + signExtendN(w>>5&0x3fff, 14)*4),
			isTbnz: isTbnz,
		}
	}
}

func (i Tbz) ObjDump(_ disasm.ViewCtx) string {
	if i.isTbnz {
		return fmt.Sprintf("tbnz %s, #0x%x, %s", i.rt, i.bit, i.target.textHex())
	}

	return fmt.Sprintf("tbz %s, #0x%x, %s", i.rt, i.bit, i.target.textHex())
}

func (i Tbz) Encode(w io.Writer, pc uint64) (int64, error) {
	target := i.target.val

	bits, err := brBits(target, int64(pc), 14)
	if err != nil {
		return 0, fmt.Errorf("tbz: %w", err)
	}

	if i.bit > 63 {
		return 0, errors.New("tbz: bit out of range")
	}

	word := uint32(0x36000000)
	if i.isTbnz {
		word = 0x37000000
	}

	if i.bit >= 32 {
		word |= 1 << 31
	}

	rt, err := armRegNum(i.rt)
	if err != nil {
		return 0, fmt.Errorf("tbz: %w", err)
	}

	return writeWord(w, word|rt|bits|i.bit&0x1f<<19)
}

func (i Tbz) MarshalJSON() ([]byte, error) {
	name := "tbz"
	if i.isTbnz {
		name = "tbnz"
	}

	return i.marshal(
		name,
		i.ObjDump(disasm.DefaultViewCtx()),
		"Branch",
		map[string]any{"Rt": i.rt, "bit": i.bit},
	)
}
