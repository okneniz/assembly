package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmcasWCtor(t *testing.T) {
	// llvm-mc-verified: amcas.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x385935cc),
		ctorWord(t, New().AmcasW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmcasW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmcasW)
	require.True(t, ok, "type = %T, want AmcasW", in)
}

func TestAmcasWDecodeEncode(t *testing.T) {
	in := decodeAmcasW(0x385935cc, 0x90000000)

	amcasw, ok := in.(AmcasW)
	require.True(t, ok, "type = %T, want AmcasW", in)
	require.Equal(t, "amcas.w $t0, $t1, $t2", amcasw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amcasw.Addr())
	require.Equal(t, 4, amcasw.Len())
	require.Equal(t, uint32(0x385935cc), ctorWord(t, amcasw))
}
