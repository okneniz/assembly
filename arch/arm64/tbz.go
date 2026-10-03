package arm64

import (
	"errors"
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Tbz — tbz rt, #bit, off (b5 selects the x/w width of Rt).
type Tbz struct {
	rt     string
	bit    uint32
	off    imm // pc-relative byte offset
	isTbnz bool
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

func (i Tbz) ObjDump(ctx disasm.ViewCtx) string {
	target := immNum(int64(ctx.Addr()) + i.off.val)
	if i.isTbnz {
		return fmt.Sprintf("tbnz %s, #0x%x, %s", i.rt, i.bit, target.textHex())
	}

	return fmt.Sprintf("tbz %s, #0x%x, %s", i.rt, i.bit, target.textHex())
}
