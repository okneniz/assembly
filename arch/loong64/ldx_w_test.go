package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxWCtor(t *testing.T) {
	// llvm-mc-verified: ldx.w $t0, $t1, $t2.
	in := New().LdxW(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x380839ac), ctorWord(t, in))

	_, ok := in.(LdxW)
	require.True(t, ok, "type = %T, want LdxW", in)
}

func TestLdxWDecodeEncode(t *testing.T) {
	in := decodeLdxW(0x380839ac, 0x90000000)

	x, ok := in.(LdxW)
	require.True(t, ok, "type = %T, want LdxW", in)
	require.Equal(t, "ldx.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x380839ac), ctorWord(t, x))
}
