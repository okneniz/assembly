package arm64

import (
	"fmt"
	"io"

	"github.com/okneniz/assembly/disasm"
)

// sysImm - system instructions with imm16 (svc/brk/hlt/hvc/udf): "#0x..", brk #0 -> "#0".
type sysImm struct {
	base

	name  string
	imm16 uint32
	enc   uint32
	shift uint // offset of the imm16 field (5 for svc/brk/udf, 21 for hlt/hvc)
}

// NewSvc — svc #imm16.
func NewSvc(imm Imm16) Instr {
	return sysImm{
		name:  "svc",
		imm16: imm.v,
		enc:   0xD4000001,
		shift: 5,
	}
}

// NewBrk — brk #imm16.
func NewBrk(imm Imm16) Instr {
	return sysImm{
		name:  "brk",
		imm16: imm.v,
		enc:   0xD4200000,
		shift: 5,
	}
}

func (i sysImm) ObjDump(_ disasm.ViewCtx) string {
	if i.name == "brk" && i.imm16 == 0 {
		return "brk #0"
	}

	return fmt.Sprintf("%s #0x%x", i.name, i.imm16)
}

func (i sysImm) Encode(w io.Writer, pc uint64) (int64, error) {
	if i.imm16 > 0xffff {
		return 0, fmt.Errorf("%s: imm out of range", i.name)
	}

	return writeWord(w, i.enc|i.imm16<<i.shift)
}

func (i sysImm) MarshalJSON() ([]byte, error) {
	return i.marshal(
		i.name,
		i.ObjDump(disasm.DefaultViewCtx()),
		"System",
		map[string]any{"imm16": i.imm16},
	)
}

func decodeSysImmOf(name string, enc uint32, shift uint) func(uint32, uint64) Instr {
	return func(w uint32, addr uint64) Instr {
		return sysImm{
			base:  newBase(addr, w),
			name:  name,
			imm16: w >> shift & 0xffff,
			enc:   enc,
			shift: shift,
		}
	}
}
