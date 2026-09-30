package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAddsImmBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		imm  int64
		sh   Sh12
		word uint32
	}{
		{"cmn x1,#0x42", "xzr", "x1", 0x42, NoSh12, 0xb101083f},
		{"adds x1,x2,#0x10,lsl#12", "x1", "x2", 0x10, LSL12, 0xb1404041},
		{"adds w2,w3,#7", "w2", "w3", 7, NoSh12, 0x31001c62},
		{"adds x2,x3,#0xfff", "x2", "x3", 0xfff, NoSh12, 0xb13ffc62},
	}
	for _, c := range cases {
		in, err := New().AddsImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm), c.sh)
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().AddsImm(reg(t, first.rd), reg(t, first.rn), imm12(t, first.imm), first.sh)
	require.NoError(t, err)
	_, ok := in.(AddsImm)
	require.True(t, ok, "type = %T, want AddsImm", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		imm  int64
		sh   Sh12
	}{
		{"adds rd sp", "sp", "x1", 1, NoSh12},
		{"adds rn xzr", "x0", "xzr", 1, NoSh12},
		{"adds x+w", "x0", "w1", 1, NoSh12},
	}
	for _, c := range errCases {
		_, err := New().AddsImm(reg(t, c.rd), reg(t, c.rn), imm12(t, c.imm), c.sh)
		assertErr(t, c.name, err)
	}
}

// TestAddsImmDecode - the Rd=31 imm word decodes as the cmn alias:
// register 31 reads as zr when S=1 (the decode used to map it to sp
// and refuse the word).
func TestAddsImmDecode(t *testing.T) {
	cases := []struct {
		word uint32
		text string
	}{
		{0xb100041f, "cmn x0, #0x1"},
		{0x3100141f, "cmn w0, #0x5"},
		{0xb140041f, "cmn x0, #0x1, lsl #12"},
	}

	for _, c := range cases {
		in, err := DecodeWord(c.word)
		require.NoError(t, err, "case %q", c.text)
		require.Equal(t, c.text, in.ObjDump(disasm.ViewCtxAt(0)), "case %q", c.text)
	}
}
