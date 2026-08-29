package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevbit8BCtor(t *testing.T) {
	// llvm-mc-verified: bitrev.8b $t0, $t1. (spec alias revbit.8b)
	require.Equal(
		t,
		uint32(0x00004dac),
		ctorWord(t, NewRevbit8B(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRevbit8B(lreg(t, 1), lreg(t, 2))
	_, ok := in.(Revbit8B)
	require.True(t, ok, "type = %T, want Revbit8B", in)
}

func TestRevbit8BDecodeEncode(t *testing.T) {
	in := decodeRevbit8B(0x00004dac, 0x90000000)

	revbit8b, ok := in.(Revbit8B)
	require.True(t, ok, "type = %T, want Revbit8B", in)
	require.Equal(t, "bitrev.8b $t0, $t1", revbit8b.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revbit8b.Addr())
	require.Equal(t, uint32(0x00004dac), ctorWord(t, revbit8b))
}
