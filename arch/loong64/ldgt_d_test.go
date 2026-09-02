package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdgtDCtor(t *testing.T) {
	// llvm-mc-verified: ldgt.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3879b9ac),
		ctorWord(t, New().LdgtD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().LdgtD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdgtD)
	require.True(t, ok, "type = %T, want LdgtD", in)
}

func TestLdgtDDecodeEncode(t *testing.T) {
	in := decodeLdgtD(0x3879b9ac, 0x90000000)

	x, ok := in.(LdgtD)
	require.True(t, ok, "type = %T, want LdgtD", in)
	require.Equal(t, "ldgt.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x3879b9ac), ctorWord(t, x))
}
