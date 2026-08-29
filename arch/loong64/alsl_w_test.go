package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAlslWCtor(t *testing.T) {
	// llvm-mc-verified: alsl.w $t0, $t1, $t2, 3 (the 1..4 shift encodes as
	// ui2 = shift - 1 = 2).
	require.Equal(
		t,
		uint32(0x000539ac),
		ctorWord(t, NewAlslW(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))),
	)

	in := NewAlslW(lreg(t, 1), lreg(t, 2), lreg(t, 3), shift3v(t, 1))
	_, ok := in.(AlslW)
	require.True(t, ok, "type = %T, want AlslW", in)
}

func TestAlslWDecodeEncode(t *testing.T) {
	in := decodeAlslW(0x000539ac, 0x90000000)

	x, ok := in.(AlslW)
	require.True(t, ok, "type = %T, want AlslW", in)

	// The raw ui2 field is 2; the decoded shift displays field + 1 = 3.
	require.Equal(t, int64(3), x.shift.val)
	require.Equal(t, "alsl.w $t0, $t1, $t2, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x000539ac), ctorWord(t, x))

	// llvm-mc-verified: alsl.w $t0, $t1, $t2, 1 (ui2 = 0) - the shift-by-one
	// holds at the lower boundary too.
	y, ok2 := decodeAlslW(0x000439ac, 0).(AlslW)
	require.True(t, ok2, "type = %T, want AlslW", y)
	require.Equal(t, int64(1), y.shift.val)
	require.Equal(t, "alsl.w $t0, $t1, $t2, 1", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x000439ac), ctorWord(t, y))
}
