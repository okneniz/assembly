package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxHuCtor(t *testing.T) {
	// llvm-mc-verified: ldx.hu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x382439ac),
		ctorWord(t, New().LdxHu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)
}

func TestLdxHuDecodeEncode(t *testing.T) {
	in := decodeLdxHu(0x382439ac, 0x90000000)

	x, ok := in.(LdxHu)
	require.True(t, ok, "type = %T, want LdxHu", in)
	require.Equal(t, "ldx.hu $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x382439ac), ctorWord(t, x))
}
