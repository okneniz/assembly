package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevb2HCtor(t *testing.T) {
	// llvm-mc-verified: revb.2h $t0, $t1.
	require.Equal(
		t,
		uint32(0x000031ac),
		ctorWord(t, New().Revb2H(lreg(t, 12), lreg(t, 13))),
	)

	in := New().Revb2H(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Revb2H)
	require.True(t, ok, "type = %T, want Revb2H", in)
}

func TestRevb2HDecodeEncode(t *testing.T) {
	in := decodeRevb2H(0x000031ac, 0x90000000)

	revb2h, ok := in.(Revb2H)
	require.True(t, ok, "type = %T, want Revb2H", in)
	require.Equal(t, "revb.2h $t0, $t1", revb2h.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revb2h.Addr())
	require.Equal(t, uint32(0x000031ac), ctorWord(t, revb2h))
}
