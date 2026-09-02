package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRevbitDCtor(t *testing.T) {
	// llvm-mc-verified: bitrev.d $t0, $t1. (spec alias revbit.d)
	require.Equal(
		t,
		uint32(0x000055ac),
		ctorWord(t, New().RevbitD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().RevbitD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RevbitD)
	require.True(t, ok, "type = %T, want RevbitD", in)
}

func TestRevbitDDecodeEncode(t *testing.T) {
	in := decodeRevbitD(0x000055ac, 0x90000000)

	revbitd, ok := in.(RevbitD)
	require.True(t, ok, "type = %T, want RevbitD", in)
	require.Equal(t, "bitrev.d $t0, $t1", revbitd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), revbitd.Addr())
	require.Equal(t, uint32(0x000055ac), ctorWord(t, revbitd))
}
