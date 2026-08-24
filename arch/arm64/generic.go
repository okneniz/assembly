package arm64

import (
	"fmt"
	"io"
	"strings"

	"github.com/okneniz/assembly/disasm"
)

// Generic — a word recognized by the generated armISA tail (step B): the
// encoding is known from the official XML (mnemonic + field offsets), but
// there is no handwritten schema/formatter for it. Printed honestly —
// mnemonic + raw fields; objdump syntax is unreachable without a formatter.
// Assembly is impossible: the instruction is decode-only (Encode returns
// an error), the assembler does not parse such text.
type Generic struct {
	base

	name   string
	fields []Field
	word   uint32
}

func decodeGeneric(e *armISAEntry, w uint32, addr uint64) Instr {
	return Generic{
		base:   newBase(addr, w),
		name:   e.Name,
		fields: e.Fields,
		word:   w,
	}
}

// fieldValue — the raw value of a field in the word.
func fieldValue(w uint32, f Field) uint32 {
	return (w >> f.Offset) & ((1 << f.Width) - 1)
}

// ObjDump — "mnemonic name=0xvalue ...". Self-describing and greppable;
// nobody but a formatter knows the actual ARM operand order, and there is
// none here.
func (i Generic) ObjDump(_ disasm.ViewCtx) string {
	var b strings.Builder
	b.WriteString(i.name)
	for _, f := range i.fields {
		fmt.Fprintf(&b, " %s=0x%x", f.Name, fieldValue(i.word, f))
	}

	return b.String()
}

// Encode — an error: generic instructions are not assembled (decode-only).
// The text→bytes round trip is undefined for them: the asm grammar does
// not know this syntax, and silently emitting the raw word would be
// dishonest.
func (i Generic) Encode(w io.Writer, pc uint64) (int64, error) {
	return 0, fmt.Errorf("%s: generic instruction, assembly not supported", i.name)
}

func (i Generic) MarshalJSON() ([]byte, error) {
	fields := make(map[string]any, len(i.fields))
	for _, f := range i.fields {
		fields[f.Name] = fmt.Sprintf("0x%x", fieldValue(i.word, f))
	}

	return i.marshal(i.name, i.ObjDump(disasm.DefaultViewCtx()), "", fields)
}
