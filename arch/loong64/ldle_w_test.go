package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdleWCtor(t *testing.T) {
	// llvm-mc-verified: ldle.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387b39ac),
		ctorWord(t, New().LdleW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().LdleW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdleW)
	require.True(t, ok, "type = %T, want LdleW", in)
}

func TestLdleWDecodeEncode(t *testing.T) {
	in := decodeLdleW(0x387b39ac, 0x90000000)

	x, ok := in.(LdleW)
	require.True(t, ok, "type = %T, want LdleW", in)
	require.Equal(t, "ldle.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387b39ac), ctorWord(t, x))
}
