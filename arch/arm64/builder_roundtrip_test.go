package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

// TestBuilderDecodeRoundTripUnit — unit level: a word encoded by a
// constructor decodes back; the decoded instruction's ObjDump text
// matches the original.
func TestBuilderDecodeRoundTripUnit(t *testing.T) {
	ins := []Instr{
		buildRet(t, xreg(t, 30)),
		New().Svc(imm16(t, 0x80)),
		New().Brk(imm16(t, 0)),
		buildMovz(t, xreg(t, 9), imm16(t, 0xbeef), Hw1),
		buildMovk(t, wreg(t, 4), imm16(t, 0x1234), Hw0),
		buildAddImm(t, SP, SP, imm12(t, 0x20), NoSh12),
		buildSubImm(t, xreg(t, 0), xreg(t, 1), imm12(t, 0x42), LSL12),
		buildAddShift(t, XZR, xreg(t, 1), xreg(t, 2), imm6(t, 7), LSR),
		buildSubShift(t, wreg(t, 0), WZR, wreg(t, 2), imm6(t, 3), LSR),
		buildLdr(t, xreg(t, 0), SP, 0x40),
		buildStr(t, wreg(t, 5), xreg(t, 29), 0x14),
	}
	for _, in := range ins {
		w := buildWord(t, in)
		back := decodeOne(w, 0x1000)
		require.Equal(
			t,
			in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()),
			"%#08x",
			w,
		)
	}
}
