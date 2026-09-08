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

func decodeSysFixedOf(name, ops, group string, enc uint32) func(uint32) Instr {
	return func(w uint32) Instr {
		return sysFixed{
			base:  newBase(w),
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

func (i sysFixed) Encode(w io.Writer) (int64, error) {
	return writeWord(w, i.enc)
}

// SysFixedOf — the fixed system instruction (dmb/dsb/yield/...).
func SysFixedOf(name, ops, group string, enc uint32) sysFixed {
	return sysFixed{name: name, ops: ops, group: group, enc: enc}
}
