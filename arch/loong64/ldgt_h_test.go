package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdgtHCtor(t *testing.T) {
	// llvm-mc-verified: ldgt.h $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3878b9ac),
		ctorWord(t, New().LdgtH(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().LdgtH(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdgtH)
	require.True(t, ok, "type = %T, want LdgtH", in)
}

func TestLdgtHDecodeEncode(t *testing.T) {
	in := decodeLdgtH(0x3878b9ac, 0x90000000)

	x, ok := in.(LdgtH)
	require.True(t, ok, "type = %T, want LdgtH", in)
	require.Equal(t, "ldgt.h $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x3878b9ac), ctorWord(t, x))
}
