package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevb4HCtor(t *testing.T) {
	// llvm-mc-verified: revb.4h $t0, $t1.
	require.Equal(
		t,
		uint32(0x000035ac),
		ctorWord(t, New().Revb4H(lreg(t, 12), lreg(t, 13))),
	)

	in := New().Revb4H(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Revb4H)
	require.True(t, ok, "type = %T, want Revb4H", in)
}

func TestRevb4HDecodeEncode(t *testing.T) {
	in := decodeRevb4H(0x000035ac, 0x90000000)

	revb4h, ok := in.(Revb4H)
	require.True(t, ok, "type = %T, want Revb4H", in)
	require.Equal(t, "revb.4h $t0, $t1", revb4h.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revb4h.Addr())
	require.Equal(t, uint32(0x000035ac), ctorWord(t, revb4h))
}
