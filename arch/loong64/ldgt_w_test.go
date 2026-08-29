package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdgtWCtor(t *testing.T) {
	// llvm-mc-verified: ldgt.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387939ac),
		ctorWord(t, NewLdgtW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewLdgtW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdgtW)
	require.True(t, ok, "type = %T, want LdgtW", in)
}

func TestLdgtWDecodeEncode(t *testing.T) {
	in := decodeLdgtW(0x387939ac, 0x90000000)

	x, ok := in.(LdgtW)
	require.True(t, ok, "type = %T, want LdgtW", in)
	require.Equal(t, "ldgt.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387939ac), ctorWord(t, x))
}
