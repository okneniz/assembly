package arm64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtrBuild(t *testing.T) {
	cases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		lsb  int64
		word uint32
	}{
		{"extr x1,x2,x3,#5", "x1", "x2", "x3", 5, 0x93c31441},
		{"ror x1,x2,#4", "x1", "x2", "x2", 4, 0x93c21041},
		{"extr x1,x2,x3,#63", "x1", "x2", "x3", 63, 0x93c3fc41},
		{"extr w1,w2,w3,#5", "w1", "w2", "w3", 5, 0x13831441},
		{"ror w1,w2,#4", "w1", "w2", "w2", 4, 0x13821041},
		{"extr w31,w31,w31,#31", "wzr", "wzr", "wzr", 31, 0x139f7fff},
	}
	for _, c := range cases {
		in, err := New().Extr(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.lsb))
		require.NoError(t, err, "case %q", c.name)
		require.Equal(t, c.word, buildWord(t, in), "case %q", c.name)
	}

	first := cases[0]
	in, err := New().Extr(reg(t, first.rd), reg(t, first.rn), reg(t, first.rm), imm6(t, first.lsb))
	require.NoError(t, err)
	_, ok := in.(Extr)
	require.True(t, ok, "type = %T, want Extr", in)

	errCases := []struct {
		name string
		rd   string
		rn   string
		rm   string
		lsb  int64
	}{
		{"extr mixed w/x", "w1", "x2", "x3", 5},
		{"extr sp", "x1", "sp", "x3", 5},
		{"extr rm sp", "x1", "x2", "sp", 5},
		{"extr w lsb out of range", "w1", "w2", "w3", 32},
	}
	for _, c := range errCases {
		_, err := New().Extr(reg(t, c.rd), reg(t, c.rn), reg(t, c.rm), imm6(t, c.lsb))
		assertErr(t, c.name, err)
	}
}

// TestExtrDecode - the decoded W words of the curated tree (the X row
// existed; the W row joined it).
func TestExtrDecode(t *testing.T) {
	cases := []struct {
		word uint32
		text string
	}{
		{0x13820c20, "extr w0, w1, w2, #0x3"},
		{0x13827c5f, "ror wzr, w2, #0x1f"},
		{0x93c20c20, "extr x0, x1, x2, #0x3"},
	}
	for _, c := range cases {
		in, err := DecodeWord(c.word)
		require.NoError(t, err, "case %#x", c.word)
		require.Equal(t, c.text, in.ObjDump(nil), "case %#x", c.word)
	}
}
