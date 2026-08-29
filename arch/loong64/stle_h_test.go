package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStleHCtor(t *testing.T) {
	// llvm-mc-verified: stle.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387eb9ac),
		ctorWord(t, NewStleH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewStleH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StleH)
	require.True(t, ok, "type = %T, want StleH", in)
}

func TestStleHDecodeEncode(t *testing.T) {
	in := decodeStleH(0x387eb9ac, 0x90000000)

	x, ok := in.(StleH)
	require.True(t, ok, "type = %T, want StleH", in)
	require.Equal(t, "stle.h $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387eb9ac), ctorWord(t, x))
}
