package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestScWCtor(t *testing.T) {
	// llvm-mc-verified: sc.w $t0, $t1, 8.
	off, err := NewImm14(8)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x210009ac),
		ctorWord(t, NewScW(lreg(t, 12), lreg(t, 13), off)),
	)

	in := NewScW(lreg(t, 1), lreg(t, 2), off)
	_, ok := in.(ScW)
	require.True(t, ok, "type = %T, want ScW", in)
}

func TestScWDecodeEncode(t *testing.T) {
	in := decodeScW(0x210009ac, 0x90000000)

	scw, ok := in.(ScW)
	require.True(t, ok, "type = %T, want ScW", in)
	require.Equal(t, "sc.w $t0, $t1, 8", scw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), scw.Addr())
	require.Equal(t, 4, scw.Len())
	require.Equal(t, int64(8), scw.off.val)
	require.Equal(t, uint32(0x210009ac), ctorWord(t, scw))
}
