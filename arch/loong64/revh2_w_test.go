package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevh2WCtor(t *testing.T) {
	// llvm-mc-verified: revh.2w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000041ac),
		ctorWord(t, NewRevh2W(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRevh2W(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Revh2W)
	require.True(t, ok, "type = %T, want Revh2W", in)
}

func TestRevh2WDecodeEncode(t *testing.T) {
	in := decodeRevh2W(0x000041ac, 0x90000000)

	revh2w, ok := in.(Revh2W)
	require.True(t, ok, "type = %T, want Revh2W", in)
	require.Equal(t, "revh.2w $t0, $t1", revh2w.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revh2w.Addr())
	require.Equal(t, uint32(0x000041ac), ctorWord(t, revh2w))
}
