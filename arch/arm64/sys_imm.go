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
	shift uint // offset of the imm16 field (bits 20:5, shift 5, for the whole family)
}

func (i sysImm) ObjDump(_ disasm.ViewCtx) string {
	if i.name == "brk" && i.imm16 == 0 {
		return "brk #0"
	}

	return fmt.Sprintf("%s #0x%x", i.name, i.imm16)
}

func (i sysImm) Encode(w io.Writer) (int64, error) {
	if i.imm16 > 0xffff {
		return 0, fmt.Errorf("%s: imm out of range", i.name)
	}

	return writeWord(w, i.enc|i.imm16<<i.shift)
}

// SysImmOf — the system instruction with imm16 (svc/brk/hlt/hvc/udf).
func SysImmOf(name string, imm16 uint32, enc uint32, shift uint) Instr {
	return sysImm{
		name:  name,
		imm16: imm16,
		enc:   enc,
		shift: shift,
	}
}
