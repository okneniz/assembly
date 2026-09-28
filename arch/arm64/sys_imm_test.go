package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSvcBrkBuild(t *testing.T) {
	cases := []struct {
		name string
		imm  int64
		word uint32
	}{
		{"svc #0x80", 0x80, 0xd4001001},
	}
	for _, c := range cases {
		in := New().Svc(imm16(t, c.imm))
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	brkCases := []struct {
		name string
		imm  int64
		word uint32
	}{
		{"brk #1", 1, 0xd4200020},
		{"brk #0", 0, 0xd4200000},
	}
	for _, c := range brkCases {
		in := New().Brk(imm16(t, c.imm))
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}

	smcCases := []struct {
		name string
		imm  int64
		word uint32
	}{
		{"smc #0", 0, 0xd4000003},
		{"smc #1", 1, 0xd4000023},
	}
	for _, c := range smcCases {
		in := New().Smc(imm16(t, c.imm))
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %q", c.name)
	}
}

// TestSysImmOfBuild - the whole imm16 exception family at its one field
// position (bits 20:5): the semihosting hlt word and the PSCI call
// among them, llvm-pinned.
func TestSysImmOfBuild(t *testing.T) {
	cases := []struct {
		name string
		imm  int64
		enc  uint32
		word uint32
	}{
		{"svc", 0x80, 0xd4000001, 0xd4001001},
		{"brk", 1, 0xd4200000, 0xd4200020},
		{"hlt", 0xf000, 0xd4400000, 0xd45e0000},
		{"hlt", 1, 0xd4400000, 0xd4400020},
		{"hvc", 0, 0xd4000002, 0xd4000002},
		{"hvc", 7, 0xd4000002, 0xd40000e2},
		{"smc", 0, 0xd4000003, 0xd4000003},
	}
	for _, c := range cases {
		in := SysImmOf(c.name, uint32(c.imm), c.enc, 5)
		require.Equal(t, c.word, buildWord(t, in), "case %s #%x", c.name, c.imm)
		back, derr := decodeOne(c.word)
		_ = derr
		require.Equal(t, in.ObjDump(disasm.DefaultViewCtx()),
			back.ObjDump(disasm.DefaultViewCtx()), "case %s #%x", c.name, c.imm)
	}
}
