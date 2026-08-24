package arm64

import (
	"io"

	"github.com/okneniz/assembly/disasm"
)

// sysFixed - operandless system instructions with a fixed word (dmb/yield/dc).
type sysFixed struct {
	base

	name  string
	ops   string
	group string
	enc   uint32
}

func decodeSysFixedOf(name, ops, group string, enc uint32) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return sysFixed{
			base:  newBase(addr, w),
			name:  name,
			ops:   ops,
			group: group,
			enc:   enc,
		}
	}
}

func (i sysFixed) ObjDump(_ disasm.ViewCtx) string {
	if i.ops == "" {
		return i.name
	}

	return i.name + " " + i.ops
}

func (i sysFixed) Encode(w io.Writer, pc uint64) (int64, error) {
	return writeWord(w, i.enc)
}

func (i sysFixed) MarshalJSON() ([]byte, error) {
	return i.marshal(i.name, i.ObjDump(disasm.DefaultViewCtx()), i.group, nil)
}
