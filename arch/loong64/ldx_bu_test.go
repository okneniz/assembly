package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLdxBuCtor(t *testing.T) {
	// llvm-mc-verified: ldx.bu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x382039ac),
		ctorWord(t, NewLdxBu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)
}

func TestLdxBuDecodeEncode(t *testing.T) {
	in := decodeLdxBu(0x382039ac, 0x90000000)

	x, ok := in.(LdxBu)
	require.True(t, ok, "type = %T, want LdxBu", in)
	require.Equal(t, "ldx.bu $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x382039ac), ctorWord(t, x))
}
