package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAndnCtor(t *testing.T) {
	// llvm-mc-verified: andn $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0016b9ac),
		ctorWord(t, NewAndn(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAndn(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Andn)
	require.True(t, ok, "type = %T, want Andn", in)
}

func TestAndnDecodeEncode(t *testing.T) {
	in := decodeOne(0x0016b9ac, 0x90000000)

	x, ok := in.(Andn)
	require.True(t, ok, "type = %T, want Andn", in)
	require.Equal(t, "andn $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0016b9ac), ctorWord(t, x))
}
