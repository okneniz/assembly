package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestStleDCtor(t *testing.T) {
	// llvm-mc-verified: stle.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387fb9ac),
		ctorWord(t, NewStleD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewStleD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(StleD)
	require.True(t, ok, "type = %T, want StleD", in)
}

func TestStleDDecodeEncode(t *testing.T) {
	in := decodeStleD(0x387fb9ac, 0x90000000)

	x, ok := in.(StleD)
	require.True(t, ok, "type = %T, want StleD", in)
	require.Equal(t, "stle.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387fb9ac), ctorWord(t, x))
}
