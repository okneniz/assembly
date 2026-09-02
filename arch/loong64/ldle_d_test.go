package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdleDCtor(t *testing.T) {
	// llvm-mc-verified: ldle.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387bb9ac),
		ctorWord(t, New().LdleD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().LdleD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdleD)
	require.True(t, ok, "type = %T, want LdleD", in)
}

func TestLdleDDecodeEncode(t *testing.T) {
	in := decodeLdleD(0x387bb9ac, 0x90000000)

	x, ok := in.(LdleD)
	require.True(t, ok, "type = %T, want LdleD", in)
	require.Equal(t, "ldle.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387bb9ac), ctorWord(t, x))
}
