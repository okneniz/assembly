package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// SubExt — sub rd, rn, rm{, ext #imm3}.
type SubExt struct {
	base
	extBase
}

const (
	SubExtX uint32 = 0xCB200000
	SubExtW uint32 = 0x4B200000
)

func decodeSubExt(w uint32, addr uint64) Instr {
	return SubExt{
		base:    newBase(addr, w),
		extBase: decodeExtBase(w),
	}
}

func (i SubExt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("sub %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i SubExt) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.extWrite(w, SubExtX, SubExtW, "sub")
}

func (i SubExt) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"sub",
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
