package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tbz — tbz rt, #bit, off (b5 selects the x/w width of Rt).
type Tbz struct {
	base

	rt     string
	bit    uint32
	off    imm // pc-relative byte offset
	isTbnz bool
}

func decodeTbzOf(isTbnz bool) func(uint32) (Instr, error) {
	return func(w uint32) (Instr, error) {
		x64 := w>>31&1 == 1
		return Tbz{
			base:   newBase(w),
			rt:     armRegName(w&0x1f, x64),
			bit:    w>>19&0x1f | w>>26&0x20,
			off:    immNum(signExtendN(w>>5&0x3fff, 14) * 4),
			isTbnz: isTbnz,
		}, nil
	}
}

func (i Tbz) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.isTbnz {
		return fmt.Sprintf("tbnz %s, #0x%x, %s", i.rt, i.bit, target.textHex())
	}

	return fmt.Sprintf("tbz %s, #0x%x, %s", i.rt, i.bit, target.textHex())
}

func (i Tbz) Encode(w io.Writer) (int64, error) {
	bits, err := offBits(i.off.val, 14)
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

	return writeWord(w, word|rt|bits<<5|i.bit&0x1f<<19)
}
