package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAndCtor(t *testing.T) {
	// llvm-mc-verified: and $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0014b9ac),
		ctorWord(t, New().And(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().And(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(And)
	require.True(t, ok, "type = %T, want And", in)
}

func TestAndDecodeEncode(t *testing.T) {
	in := decodeOne(0x0014b9ac, 0x90000000)

	x, ok := in.(And)
	require.True(t, ok, "type = %T, want And", in)
	require.Equal(t, "and $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0014b9ac), ctorWord(t, x))
}
