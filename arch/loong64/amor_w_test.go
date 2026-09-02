package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmorWCtor(t *testing.T) {
	// llvm-mc-verified: amor.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386335cc),
		ctorWord(t, New().AmorW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmorW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmorW)
	require.True(t, ok, "type = %T, want AmorW", in)
}

func TestAmorWDecodeEncode(t *testing.T) {
	in := decodeAmorW(0x386335cc, 0x90000000)

	amorw, ok := in.(AmorW)
	require.True(t, ok, "type = %T, want AmorW", in)
	require.Equal(t, "amor.w $t0, $t1, $t2", amorw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amorw.Addr())
	require.Equal(t, 4, amorw.Len())
	require.Equal(t, uint32(0x386335cc), ctorWord(t, amorw))
}
