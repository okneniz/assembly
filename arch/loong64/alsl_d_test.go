package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAlslDCtor(t *testing.T) {
	// llvm-mc-verified: alsl.d $t0, $t1, $t2, 3 (ui2 = shift - 1 = 2).
	require.Equal(
		t,
		uint32(0x002d39ac),
		ctorWord(t, New().AlslD(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))),
	)

	in := New().AlslD(lreg(t, 1), lreg(t, 2), lreg(t, 3), shift3v(t, 4))
	_, ok := in.(AlslD)
	require.True(t, ok, "type = %T, want AlslD", in)
}

func TestAlslDDecodeEncode(t *testing.T) {
	in := decodeAlslD(0x002d39ac, 0x90000000)

	x, ok := in.(AlslD)
	require.True(t, ok, "type = %T, want AlslD", in)

	// The raw ui2 field is 2; the decoded shift displays field + 1 = 3.
	require.Equal(t, int64(3), x.shift.val)
	require.Equal(t, "alsl.d $t0, $t1, $t2, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x002d39ac), ctorWord(t, x))

	// llvm-mc-verified: alsl.d $t0, $t1, $t2, 4 (ui2 = 3) - the upper
	// boundary of the shift-by-one.
	y, ok2 := decodeAlslD(0x002db9ac, 0).(AlslD)
	require.True(t, ok2, "type = %T, want AlslD", y)
	require.Equal(t, int64(4), y.shift.val)
	require.Equal(t, "alsl.d $t0, $t1, $t2, 4", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x002db9ac), ctorWord(t, y))
}
