package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSraDCtor(t *testing.T) {
	// llvm-mc-verified: sra.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0019b9ac),
		ctorWord(t, NewSraD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewSraD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SraD)
	require.True(t, ok, "type = %T, want SraD", in)
}

func TestSraDDecodeEncode(t *testing.T) {
	in := decodeOne(0x0019b9ac, 0x90000000)

	x, ok := in.(SraD)
	require.True(t, ok, "type = %T, want SraD", in)
	require.Equal(t, "sra.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0019b9ac), ctorWord(t, x))
}
