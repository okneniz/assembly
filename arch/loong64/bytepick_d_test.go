package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBytepickDCtor(t *testing.T) {
	// llvm-mc-verified: bytepick.d $t0, $t1, $t2, 3 (a3 = 3 as-is, no
	// shift-by-one).
	require.Equal(
		t,
		uint32(0x000db9ac),
		ctorWord(t, New().BytepickD(lreg(t, 12), lreg(t, 13), lreg(t, 14), uimm3v(t, 3))),
	)

	in := New().BytepickD(lreg(t, 1), lreg(t, 2), lreg(t, 3), uimm3v(t, 7))
	_, ok := in.(BytepickD)
	require.True(t, ok, "type = %T, want BytepickD", in)
}

func TestBytepickDDecodeEncode(t *testing.T) {
	in := decodeBytepickD(0x000db9ac, 0x90000000)

	x, ok := in.(BytepickD)
	require.True(t, ok, "type = %T, want BytepickD", in)
	require.Equal(t, int64(3), x.sel.val)
	require.Equal(t, "bytepick.d $t0, $t1, $t2, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x000db9ac), ctorWord(t, x))

	// llvm-mc-verified: bytepick.d $t0, $t1, $t2, 7 - the full a3 field.
	y, ok2 := decodeBytepickD(0x000fb9ac, 0).(BytepickD)
	require.True(t, ok2, "type = %T, want BytepickD", y)
	require.Equal(t, "bytepick.d $t0, $t1, $t2, 7", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x000fb9ac), ctorWord(t, y))
}
