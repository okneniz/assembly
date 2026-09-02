package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdpteCtor(t *testing.T) {
	v, err := New().UImm8(1)
	require.NoError(t, err)

	// llvm-mc-verified: ldpte $t1, 1.
	require.Equal(
		t,
		uint32(0x064405a0),
		ctorWord(t, New().Ldpte(lreg(t, 13), v)),
	)

	in := New().Ldpte(lreg(t, 13), v)
	_, ok := in.(Ldpte)
	require.True(t, ok, "type = %T, want Ldpte", in)
}

func TestLdpteDecodeEncode(t *testing.T) {
	// llvm-mc-verified: ldpte $t1, 1.
	in := decodeLdpte(0x064405a0, 0x90000000)

	x, ok := in.(Ldpte)
	require.True(t, ok, "type = %T, want Ldpte", in)
	require.Equal(t, "ldpte $t1, 1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064405a0), ctorWord(t, x))
}
