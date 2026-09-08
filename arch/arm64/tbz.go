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

// Tbz — tbz rt, #bit, off: off — the pc-relative byte offset of the
// branch destination (the ±32KB imm14 range is checked at encode time;
// the absolute target is off + the instruction address). The register
// width is dictated by the bit number — the sf bit
// of the encoding is the bit's b5: bits 32..63 need an x register, bits
// 0..31 — a w one (register 31 reads as zr — use XZR/WZR).
func (Builder) Tbz(rt Reg, bit uint32, off int64) (Instr, error) {
	if err := requireClass(rt, "Tbz", "rt", "x/w register (register 31 reads as zr — use XZR/WZR)",
		classX, classW, classXZR, classWZR); err != nil {
		return nil, err
	}

	if bit > 63 {
		return nil, fmt.Errorf("arm64.NewTbz: operand bit: %d is out of 0..63", bit)
	}

	if rt.Is64() != (bit >= 32) {
		want := "w"
		if bit >= 32 {
			want = "x"
		}

		return nil, fmt.Errorf(
			"arm64.NewTbz: operand bit: %d needs a %s register (the sf bit of the encoding is the bit's b5), got %s",
			bit,
			want,
			rt.name(),
		)
	}

	return Tbz{rt: rt.name(), bit: bit, off: immNum(off), isTbnz: false}, nil
}

func decodeTbzOf(isTbnz bool) func(uint32) Instr {
	return func(w uint32) Instr {
		x64 := w>>31&1 == 1
		return Tbz{
			base:   newBase(w),
			rt:     armRegName(w&0x1f, x64),
			bit:    w>>19&0x1f | w>>26&0x20,
			off:    immNum(signExtendN(w>>5&0x3fff, 14) * 4),
			isTbnz: isTbnz,
		}
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
