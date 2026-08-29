package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrccWWWCtor(t *testing.T) {
	// llvm-mc-verified: crcc.w.w.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002739ac),
		ctorWord(t, NewCrccWWW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewCrccWWW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrccWWW)
	require.True(t, ok, "type = %T, want CrccWWW", in)
}

func TestCrccWWWDecodeEncode(t *testing.T) {
	in := decodeCrccWWW(0x002739ac, 0x90000000)

	x, ok := in.(CrccWWW)
	require.True(t, ok, "type = %T, want CrccWWW", in)
	require.Equal(t, "crcc.w.w.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x002739ac), ctorWord(t, x))
}
