package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminWuCtor(t *testing.T) {
	// llvm-mc-verified: ammin.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386835cc),
		ctorWord(t, NewAmminWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmminWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminWu)
	require.True(t, ok, "type = %T, want AmminWu", in)
}

func TestAmminWuDecodeEncode(t *testing.T) {
	in := decodeAmminWu(0x386835cc, 0x90000000)

	amminwu, ok := in.(AmminWu)
	require.True(t, ok, "type = %T, want AmminWu", in)
	require.Equal(t, "ammin.wu $t0, $t1, $t2", amminwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amminwu.Addr())
	require.Equal(t, 4, amminwu.Len())
	require.Equal(t, uint32(0x386835cc), ctorWord(t, amminwu))
}
