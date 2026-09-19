package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// Dup — dup.Arr vd, wn (DUP general: every lane = the integer
// register's value; .2d takes an x register, the rest a w one).
type Dup struct {
	base

	vd, gpr string
	arr     string // 8b/16b/4h/8h/2s/4s/2d
}

// newDup - the Dup constructor: validates the operands and assembles
// the struct (the Builder method delegates here; the decoder calls it
// with values read from the word).
func newDup(b base, vd VReg, gpr Reg, arr string) (Dup, error) {
	err := requireArr("Dup", arr, "8b", "16b", "4h", "8h", "2s", "4s", "2d")
	if err != nil {
		return Dup{}, err
	}

	err = requireGprClass(gpr, "Dup", "wn")
	if err != nil {
		return Dup{}, err
	}

	_, size, err := arrBits(arr)
	if err != nil {
		return Dup{}, err
	}

	if (size == 3) != gpr.Is64() {
		return Dup{}, fmt.Errorf(
			"arm64.NewDup: register widths must match: %s vs %s (.2d takes an x register)",
			arr, gpr.name(),
		)
	}

	return Dup{
		base: b,
		vd:   vd.name(),
		gpr:  gpr.name(),
		arr:  arr,
	}, nil
}

const dupEnc uint32 = 0x0E000000 // dup vd, wn (family base; op bits 0)

func (i Dup) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("dup.%s %s, %s", i.arr, i.vd, i.gpr)
}

func (i Dup) Encode(w io.Writer) (int64, error) {
	vd, err := armRegNum(i.vd)
	if err != nil {
		return 0, fmt.Errorf("dup: %w", err)
	}

	gpr, err := armRegNum(i.gpr)
	if err != nil {
		return 0, fmt.Errorf("dup: %w", err)
	}

	q, size, err := arrBits(i.arr)
	if err != nil {
		return 0, fmt.Errorf("dup: %w", err)
	}

	return writeWord(w, dupEnc|q<<30|0xC00|1<<size<<16|gpr<<5|vd)
}

func (Builder) Dup(vd VReg, wn Reg, arr string) (Instr, error) {
	return newDup(base{}, vd, wn, arr)
}
