package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestModWCtor(t *testing.T) {
	// llvm-mc-verified: mod.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0020b9ac),
		ctorWord(t, New().ModW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().ModW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(ModW)
	require.True(t, ok, "type = %T, want ModW", in)
}

func TestModWDecodeEncode(t *testing.T) {
	in := decodeModW(0x0020b9ac, 0x90000000)

	modw, ok := in.(ModW)
	require.True(t, ok, "type = %T, want ModW", in)
	require.Equal(t, "mod.w $t0, $t1, $t2", modw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), modw.Addr())
	require.Equal(t, uint32(0x0020b9ac), ctorWord(t, modw))
}
