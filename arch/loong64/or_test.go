package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestOrCtor(t *testing.T) {
	// llvm-mc-verified: or $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001539ac),
		ctorWord(t, NewOr(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewOr(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Or)
	require.True(t, ok, "type = %T, want Or", in)
}

func TestOrDecodeEncode(t *testing.T) {
	in := decodeOne(0x001539ac, 0x90000000)

	x, ok := in.(Or)
	require.True(t, ok, "type = %T, want Or", in)
	require.Equal(t, "or $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001539ac), ctorWord(t, x))
}
