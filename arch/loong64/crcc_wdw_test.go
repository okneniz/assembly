package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrccWDWCtor(t *testing.T) {
	// llvm-mc-verified: crcc.w.d.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0027b9ac),
		ctorWord(t, New().CrccWDW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().CrccWDW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrccWDW)
	require.True(t, ok, "type = %T, want CrccWDW", in)
}

func TestCrccWDWDecodeEncode(t *testing.T) {
	in := decodeCrccWDW(0x0027b9ac, 0x90000000)

	x, ok := in.(CrccWDW)
	require.True(t, ok, "type = %T, want CrccWDW", in)
	require.Equal(t, "crcc.w.d.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0027b9ac), ctorWord(t, x))
}
