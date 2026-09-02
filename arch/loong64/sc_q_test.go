package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestScQCtor(t *testing.T) {
	// llvm-mc-verified: sc.q $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385735cc),
		ctorWord(t, New().ScQ(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().ScQ(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(ScQ)
	require.True(t, ok, "type = %T, want ScQ", in)
}

func TestScQDecodeEncode(t *testing.T) {
	in := decodeScQ(0x385735cc, 0x90000000)

	scq, ok := in.(ScQ)
	require.True(t, ok, "type = %T, want ScQ", in)
	require.Equal(t, "sc.q $t0, $t1, $t2", scq.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), scq.Addr())
	require.Equal(t, 4, scq.Len())
	require.Equal(t, uint32(0x385735cc), ctorWord(t, scq))
}
