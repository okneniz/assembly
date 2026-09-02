package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdgtBCtor(t *testing.T) {
	// llvm-mc-verified: ldgt.b $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x387839ac),
		ctorWord(t, New().LdgtB(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().LdgtB(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(LdgtB)
	require.True(t, ok, "type = %T, want LdgtB", in)
}

func TestLdgtBDecodeEncode(t *testing.T) {
	in := decodeLdgtB(0x387839ac, 0x90000000)

	x, ok := in.(LdgtB)
	require.True(t, ok, "type = %T, want LdgtB", in)
	require.Equal(t, "ldgt.b $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x387839ac), ctorWord(t, x))
}
