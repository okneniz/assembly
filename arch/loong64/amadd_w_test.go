package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmaddWCtor(t *testing.T) {
	// llvm-mc-verified: amadd.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386135cc),
		ctorWord(t, NewAmaddW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmaddW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmaddW)
	require.True(t, ok, "type = %T, want AmaddW", in)
}

func TestAmaddWDecodeEncode(t *testing.T) {
	in := decodeAmaddW(0x386135cc, 0x90000000)

	amaddw, ok := in.(AmaddW)
	require.True(t, ok, "type = %T, want AmaddW", in)
	require.Equal(t, "amadd.w $t0, $t1, $t2", amaddw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amaddw.Addr())
	require.Equal(t, 4, amaddw.Len())
	require.Equal(t, uint32(0x386135cc), ctorWord(t, amaddw))
}
