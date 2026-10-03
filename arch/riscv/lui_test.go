package riscv

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	parsec "github.com/okneniz/parsec"
	parsecbytes "github.com/okneniz/parsec/bytes"
)

func TestLuiCtor(t *testing.T) {
	require.Equal(
		t,
		uint32(0x123452b7),
		ctorWord(t, New().Lui(xreg(t, 5), imm20(t, 0x12345))),
	)
	in := New().Lui(xreg(t, 1), imm20(t, 1))
	_, ok := in.(Lui)
	require.True(t, ok, "type = %T, want Lui", in)
	// Imm20 outside 0..0xfffff is rejected by the Imm20 constructor.
	_, err := New().Imm20(-1)
	require.Error(t, err)
	_, err = New().Imm20(0x100000)
	require.Error(t, err)
}

// TestCLuiCanon - the decoded c.lui keeps the RAW UNSIGNED imm20 in
// the imm slot (llvm-mc prints "lui a0, 1" for the halfword 0x6505 and
// the unsigned form for the negative nzimm window): the text
// re-assembles back to the same compressed bytes.
func TestCLuiCanon(t *testing.T) {
	for _, c := range []struct {
		half [2]byte
		text string
	}{
		{[2]byte{0x05, 0x65}, "lui a0, 0x1"},
		{[2]byte{0xb9, 0x70}, "lui ra, 0xfffee"}, // c.lui ra, -18
	} {
		in, err := MakeDecoder()(
			parsec.Stateless{},
			parsecbytes.Buffer(c.half[:]),
		)
		require.NoError(t, err)
		require.Len(t, in, 1)
		require.Equal(t, c.text, in[0].ObjDump(nil))

		var buf bytes.Buffer
		_, eerr := in[0].Encode(&buf, EncOpts{})
		require.NoError(t, eerr)
		require.Equal(t, c.half[:], buf.Bytes())
	}
}
