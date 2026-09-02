package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBstrinsDCtor(t *testing.T) {
	// llvm-mc-verified: bstrins.d $t0, $t1, 5, 3 (msb at [21:16], lsb at
	// [15:10]).
	require.Equal(
		t,
		uint32(0x00850dac),
		ctorWord(t, New().BstrinsD(lreg(t, 12), lreg(t, 13), uimm6v(t, 5), uimm6v(t, 3))),
	)

	in := New().BstrinsD(lreg(t, 1), lreg(t, 2), uimm6v(t, 63), uimm6v(t, 0))
	_, ok := in.(BstrinsD)
	require.True(t, ok, "type = %T, want BstrinsD", in)
}

func TestBstrinsDDecodeEncode(t *testing.T) {
	in := decodeBstrinsD(0x00850dac, 0x90000000)

	x, ok := in.(BstrinsD)
	require.True(t, ok, "type = %T, want BstrinsD", in)
	require.Equal(t, int64(5), x.msb.val)
	require.Equal(t, int64(3), x.lsb.val)
	require.Equal(t, "bstrins.d $t0, $t1, 5, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x00850dac), ctorWord(t, x))

	// llvm-mc-verified: bstrins.d $t0, $t1, 63, 0 - the full 6-bit fields.
	y, ok2 := decodeBstrinsD(0x00bf01ac, 0).(BstrinsD)
	require.True(t, ok2, "type = %T, want BstrinsD", y)
	require.Equal(t, "bstrins.d $t0, $t1, 63, 0", y.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x00bf01ac), ctorWord(t, y))
}
