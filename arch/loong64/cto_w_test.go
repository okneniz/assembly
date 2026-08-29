package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCtoWCtor(t *testing.T) {
	// llvm-mc-verified: cto.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x000019ac),
		ctorWord(t, NewCtoW(lreg(t, 12), lreg(t, 13))),
	)

	in := NewCtoW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(CtoW)
	require.True(t, ok, "type = %T, want CtoW", in)
}

func TestCtoWDecodeEncode(t *testing.T) {
	in := decodeCtoW(0x000019ac, 0x90000000)

	ctow, ok := in.(CtoW)
	require.True(t, ok, "type = %T, want CtoW", in)
	require.Equal(t, "cto.w $t0, $t1", ctow.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ctow.Addr())
	require.Equal(t, uint32(0x000019ac), ctorWord(t, ctow))
}
