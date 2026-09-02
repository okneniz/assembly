package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRdtimelWCtor(t *testing.T) {
	// llvm-mc-verified: rdtimel.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000061ac),
		ctorWord(t, New().RdtimelW(lreg(t, 12), lreg(t, 13))),
	)

	in := New().RdtimelW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RdtimelW)
	require.True(t, ok, "type = %T, want RdtimelW", in)
}

func TestRdtimelWDecodeEncode(t *testing.T) {
	in := decodeRdtimelW(0x000061ac, 0x90000000)

	rdtimelw, ok := in.(RdtimelW)
	require.True(t, ok, "type = %T, want RdtimelW", in)
	require.Equal(t, "rdtimel.w $t0, $t1", rdtimelw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), rdtimelw.Addr())
	require.Equal(t, uint32(0x000061ac), ctorWord(t, rdtimelw))
}
