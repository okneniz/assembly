package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevbitWCtor(t *testing.T) {
	// llvm-mc-verified: bitrev.w $t0, $t1. (spec alias revbit.w)
	require.Equal(
		t,
		uint32(0x000051ac),
		ctorWord(t, NewRevbitW(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRevbitW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RevbitW)
	require.True(t, ok, "type = %T, want RevbitW", in)
}

func TestRevbitWDecodeEncode(t *testing.T) {
	in := decodeRevbitW(0x000051ac, 0x90000000)

	revbitw, ok := in.(RevbitW)
	require.True(t, ok, "type = %T, want RevbitW", in)
	require.Equal(t, "bitrev.w $t0, $t1", revbitw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revbitw.Addr())
	require.Equal(t, uint32(0x000051ac), ctorWord(t, revbitw))
}
