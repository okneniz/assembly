package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevhDCtor(t *testing.T) {
	// llvm-mc-verified: revh.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x000045ac),
		ctorWord(t, NewRevhD(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRevhD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RevhD)
	require.True(t, ok, "type = %T, want RevhD", in)
}

func TestRevhDDecodeEncode(t *testing.T) {
	in := decodeRevhD(0x000045ac, 0x90000000)

	revhd, ok := in.(RevhD)
	require.True(t, ok, "type = %T, want RevhD", in)
	require.Equal(t, "revh.d $t0, $t1", revhd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revhd.Addr())
	require.Equal(t, uint32(0x000045ac), ctorWord(t, revhd))
}
