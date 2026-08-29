package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSraWCtor(t *testing.T) {
	// llvm-mc-verified: sra.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001839ac),
		ctorWord(t, NewSraW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewSraW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SraW)
	require.True(t, ok, "type = %T, want SraW", in)
}

func TestSraWDecodeEncode(t *testing.T) {
	in := decodeOne(0x001839ac, 0x90000000)

	x, ok := in.(SraW)
	require.True(t, ok, "type = %T, want SraW", in)
	require.Equal(t, "sra.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001839ac), ctorWord(t, x))
}
