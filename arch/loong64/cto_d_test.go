package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCtoDCtor(t *testing.T) {
	// llvm-mc-verified: cto.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x000029ac),
		ctorWord(t, New().CtoD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().CtoD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(CtoD)
	require.True(t, ok, "type = %T, want CtoD", in)
}

func TestCtoDDecodeEncode(t *testing.T) {
	in := decodeCtoD(0x000029ac, 0x90000000)

	ctod, ok := in.(CtoD)
	require.True(t, ok, "type = %T, want CtoD", in)
	require.Equal(t, "cto.d $t0, $t1", ctod.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ctod.Addr())
	require.Equal(t, uint32(0x000029ac), ctorWord(t, ctod))
}
