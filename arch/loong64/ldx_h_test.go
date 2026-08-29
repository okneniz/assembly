package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxHCtor(t *testing.T) {
	// llvm-mc-verified: ldx.h $t0, $t1, $t2.
	in := NewLdxH(lreg(t, 12), lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x380439ac), ctorWord(t, in))

	_, ok := in.(LdxH)
	require.True(t, ok, "type = %T, want LdxH", in)
}

func TestLdxHDecodeEncode(t *testing.T) {
	in := decodeLdxH(0x380439ac, 0x90000000)

	x, ok := in.(LdxH)
	require.True(t, ok, "type = %T, want LdxH", in)
	require.Equal(t, "ldx.h $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x380439ac), ctorWord(t, x))
}
