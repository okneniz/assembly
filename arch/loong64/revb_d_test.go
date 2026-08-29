package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevbDCtor(t *testing.T) {
	// llvm-mc-verified: revb.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x00003dac),
		ctorWord(t, NewRevbD(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRevbD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RevbD)
	require.True(t, ok, "type = %T, want RevbD", in)
}

func TestRevbDDecodeEncode(t *testing.T) {
	in := decodeRevbD(0x00003dac, 0x90000000)

	revbd, ok := in.(RevbD)
	require.True(t, ok, "type = %T, want RevbD", in)
	require.Equal(t, "revb.d $t0, $t1", revbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revbd.Addr())
	require.Equal(t, uint32(0x00003dac), ctorWord(t, revbd))
}
