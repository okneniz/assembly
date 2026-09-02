package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrccWBWCtor(t *testing.T) {
	// llvm-mc-verified: crcc.w.b.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002639ac),
		ctorWord(t, New().CrccWBW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().CrccWBW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrccWBW)
	require.True(t, ok, "type = %T, want CrccWBW", in)
}

func TestCrccWBWDecodeEncode(t *testing.T) {
	in := decodeCrccWBW(0x002639ac, 0x90000000)

	x, ok := in.(CrccWBW)
	require.True(t, ok, "type = %T, want CrccWBW", in)
	require.Equal(t, "crcc.w.b.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x002639ac), ctorWord(t, x))
}
