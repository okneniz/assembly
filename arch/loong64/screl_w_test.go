package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestScrelWCtor(t *testing.T) {
	// llvm-mc-verified: screl.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x385785ac),
		ctorWord(t, New().ScrelW(lreg(t, 12), lreg(t, 13))),
	)

	in := New().ScrelW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(ScrelW)
	require.True(t, ok, "type = %T, want ScrelW", in)
}

func TestScrelWDecodeEncode(t *testing.T) {
	in := decodeScrelW(0x385785ac, 0x90000000)

	screlw, ok := in.(ScrelW)
	require.True(t, ok, "type = %T, want ScrelW", in)
	require.Equal(t, "screl.w $t0, $t1", screlw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), screlw.Addr())
	require.Equal(t, 4, screlw.Len())
	require.Equal(t, uint32(0x385785ac), ctorWord(t, screlw))
}
