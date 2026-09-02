package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCloWCtor(t *testing.T) {
	// llvm-mc-verified: clo.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000011ac),
		ctorWord(t, New().CloW(lreg(t, 12), lreg(t, 13))),
	)

	in := New().CloW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(CloW)
	require.True(t, ok, "type = %T, want CloW", in)
}

func TestCloWDecodeEncode(t *testing.T) {
	in := decodeCloW(0x000011ac, 0x90000000)

	clow, ok := in.(CloW)
	require.True(t, ok, "type = %T, want CloW", in)
	require.Equal(t, "clo.w $t0, $t1", clow.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), clow.Addr())
	require.Equal(t, uint32(0x000011ac), ctorWord(t, clow))
}
