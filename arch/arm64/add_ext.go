package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// AddExt — add rd, rn, rm{, ext #imm3}.
type AddExt struct {
	base
	extBase
}

const (
	AddExtX uint32 = 0x8B200000
	AddExtW uint32 = 0x0B200000
)

func decodeAddExt(w uint32, addr uint64) Instr {
	return AddExt{
		base:    newBase(addr, w),
		extBase: decodeExtBase(w),
	}
}

func (i AddExt) ObjDump(_ disasm.ViewCtx) string {
	return fmt.Sprintf("add %s, %s, %s%s", addSubRegName(i.rdNum, i.isf, false),
		addSubRegName(i.rnNum, i.isf, false), addSubRegName(i.rmNum, i.isf, false), i.extMod(false))
}

func (i AddExt) Encode(w io.Writer, pc uint64) (int64, error) {
	return i.extWrite(w, AddExtX, AddExtW, "add")
}

func (i AddExt) MarshalJSON() ([]byte, error) {
	return i.marshal(
		"add",
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
