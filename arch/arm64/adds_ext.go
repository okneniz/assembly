package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddsExt — adds rd, rn, rm{, ext #imm3}; pseudo: cmn (Rd=31).
type AddsExt struct {
	base
	extBase
}

const (
	AddsExtX uint32 = 0xAB200000
	AddsExtW uint32 = 0x2B200000
)

func decodeAddsExt(w uint32, addr uint64) Instr {
	return AddsExt{
		base:    newBase(addr, w),
		extBase: decodeExtBase(w),
	}
}

func (i AddsExt) ObjDump(_ disasm.ViewCtx) string {
	if i.rdNum == 31 {
		rnz := addSubRegName(i.rnNum, i.isf, false)
		return fmt.Sprintf(
			"cmn %s, %s%s",
			rnz,
			addSubRegName(i.rmNum, i.isf, false),
			i.extMod(true),
		)
	}

	return fmt.Sprintf("adds %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i AddsExt) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.extWrite(w, AddsExtX, AddsExtW, "adds")
}

func (i AddsExt) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"adds",
		i.ObjDump(disasm.DefaultViewCtx()),
		"Data processing - register",
		map[string]any{
			"Rd":     i.rdNum,
			"Rn":     i.rnNum,
			"Rm":     i.rmNum,
			"option": i.option,
			"imm3":   i.imm3,
		},
	)
}
