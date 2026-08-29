package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCtzWCtor(t *testing.T) {
	// llvm-mc-verified: ctz.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x00001dac),
		ctorWord(t, NewCtzW(lreg(t, 12), lreg(t, 13))),
	)

	in := NewCtzW(lreg(t, 1), lreg(t, 2))
	_, ok := in.(CtzW)
	require.True(t, ok, "type = %T, want CtzW", in)
}

func TestCtzWDecodeEncode(t *testing.T) {
	in := decodeCtzW(0x00001dac, 0x90000000)

	ctzw, ok := in.(CtzW)
	require.True(t, ok, "type = %T, want CtzW", in)
	require.Equal(t, "ctz.w $t0, $t1", ctzw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ctzw.Addr())
	require.Equal(t, uint32(0x00001dac), ctorWord(t, ctzw))
}
