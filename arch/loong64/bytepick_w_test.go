package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBytepickWCtor(t *testing.T) {
	// llvm-mc-verified: bytepick.w $t0, $t1, $t2, 3 (a2 = 3 as-is, no
	// shift-by-one).
	require.Equal(
		t,
		uint32(0x0009b9ac),
		ctorWord(t, New().BytepickW(lreg(t, 12), lreg(t, 13), lreg(t, 14), uimm2v(t, 3))),
	)

	in := New().BytepickW(lreg(t, 1), lreg(t, 2), lreg(t, 3), uimm2v(t, 0))
	_, ok := in.(BytepickW)
	require.True(t, ok, "type = %T, want BytepickW", in)
}

func TestBytepickWDecodeEncode(t *testing.T) {
	in := decodeBytepickW(0x0009b9ac, 0x90000000)

	x, ok := in.(BytepickW)
	require.True(t, ok, "type = %T, want BytepickW", in)
	require.Equal(t, int64(3), x.sel.val)
	require.Equal(t, "bytepick.w $t0, $t1, $t2, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0009b9ac), ctorWord(t, x))
}
