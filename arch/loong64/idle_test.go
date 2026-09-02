package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIdleCtor(t *testing.T) {
	code, err := New().Code15(0)
	require.NoError(t, err)

	// llvm-mc-verified: idle 0.
	require.Equal(
		t,
		uint32(0x06488000),
		ctorWord(t, New().Idle(code)),
	)

	in := New().Idle(code)
	_, ok := in.(Idle)
	require.True(t, ok, "type = %T, want Idle", in)
}

func TestIdleDecodeEncode(t *testing.T) {
	// llvm-mc-verified: idle 1.
	in := decodeIdle(0x06488001, 0x90000000)

	x, ok := in.(Idle)
	require.True(t, ok, "type = %T, want Idle", in)
	require.Equal(t, "idle 1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06488001), ctorWord(t, x))
}
