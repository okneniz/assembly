package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCrccWHWCtor(t *testing.T) {
	// llvm-mc-verified: crcc.w.h.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0026b9ac),
		ctorWord(t, New().CrccWHW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().CrccWHW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(CrccWHW)
	require.True(t, ok, "type = %T, want CrccWHW", in)
}

func TestCrccWHWDecodeEncode(t *testing.T) {
	in := decodeCrccWHW(0x0026b9ac, 0x90000000)

	x, ok := in.(CrccWHW)
	require.True(t, ok, "type = %T, want CrccWHW", in)
	require.Equal(t, "crcc.w.h.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0026b9ac), ctorWord(t, x))
}
