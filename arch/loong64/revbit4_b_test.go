package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevbit4BCtor(t *testing.T) {
	// llvm-mc-verified: bitrev.4b $t0, $t1. (spec alias revbit.4b)
	require.Equal(
		t,
		uint32(0x000049ac),
		ctorWord(t, New().Revbit4B(lreg(t, 12), lreg(t, 13))),
	)

	in := New().Revbit4B(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Revbit4B)
	require.True(t, ok, "type = %T, want Revbit4B", in)
}

func TestRevbit4BDecodeEncode(t *testing.T) {
	in := decodeRevbit4B(0x000049ac, 0x90000000)

	revbit4b, ok := in.(Revbit4B)
	require.True(t, ok, "type = %T, want Revbit4B", in)
	require.Equal(t, "bitrev.4b $t0, $t1", revbit4b.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revbit4b.Addr())
	require.Equal(t, uint32(0x000049ac), ctorWord(t, revbit4b))
}
