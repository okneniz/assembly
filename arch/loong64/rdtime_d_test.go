package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestRdtimeDCtor(t *testing.T) {
	// llvm-mc-verified: rdtime.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x000069ac),
		ctorWord(t, NewRdtimeD(lreg(t, 12), lreg(t, 13))),
	)

	in := NewRdtimeD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(RdtimeD)
	require.True(t, ok, "type = %T, want RdtimeD", in)
}

func TestRdtimeDDecodeEncode(t *testing.T) {
	in := decodeRdtimeD(0x000069ac, 0x90000000)

	rdtimed, ok := in.(RdtimeD)
	require.True(t, ok, "type = %T, want RdtimeD", in)
	require.Equal(t, "rdtime.d $t0, $t1", rdtimed.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), rdtimed.Addr())
	require.Equal(t, uint32(0x000069ac), ctorWord(t, rdtimed))
}
