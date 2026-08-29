package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminWCtor(t *testing.T) {
	// llvm-mc-verified: ammin.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386635cc),
		ctorWord(t, NewAmminW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmminW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminW)
	require.True(t, ok, "type = %T, want AmminW", in)
}

func TestAmminWDecodeEncode(t *testing.T) {
	in := decodeAmminW(0x386635cc, 0x90000000)

	amminw, ok := in.(AmminW)
	require.True(t, ok, "type = %T, want AmminW", in)
	require.Equal(t, "ammin.w $t0, $t1, $t2", amminw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amminw.Addr())
	require.Equal(t, 4, amminw.Len())
	require.Equal(t, uint32(0x386635cc), ctorWord(t, amminw))
}
