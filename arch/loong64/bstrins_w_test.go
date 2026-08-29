package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestBstrinsWCtor(t *testing.T) {
	// llvm-mc-verified: bstrins.w $t0, $t1, 5, 3 (msb at [20:16], lsb at
	// [14:10]).
	require.Equal(
		t,
		uint32(0x00650dac),
		ctorWord(t, NewBstrinsW(lreg(t, 12), lreg(t, 13), uimm5v(t, 5), uimm5v(t, 3))),
	)

	in := NewBstrinsW(lreg(t, 1), lreg(t, 2), uimm5v(t, 31), uimm5v(t, 0))
	_, ok := in.(BstrinsW)
	require.True(t, ok, "type = %T, want BstrinsW", in)
}

func TestBstrinsWDecodeEncode(t *testing.T) {
	in := decodeBstrinsW(0x00650dac, 0x90000000)

	x, ok := in.(BstrinsW)
	require.True(t, ok, "type = %T, want BstrinsW", in)
	require.Equal(t, int64(5), x.msb.val)
	require.Equal(t, int64(3), x.lsb.val)
	require.Equal(t, "bstrins.w $t0, $t1, 5, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x00650dac), ctorWord(t, x))
}
