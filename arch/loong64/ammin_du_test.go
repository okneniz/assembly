package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmminDuCtor(t *testing.T) {
	// llvm-mc-verified: ammin.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x3868b5cc),
		ctorWord(t, NewAmminDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmminDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmminDu)
	require.True(t, ok, "type = %T, want AmminDu", in)
}

func TestAmminDuDecodeEncode(t *testing.T) {
	in := decodeAmminDu(0x3868b5cc, 0x90000000)

	ammindu, ok := in.(AmminDu)
	require.True(t, ok, "type = %T, want AmminDu", in)
	require.Equal(t, "ammin.du $t0, $t1, $t2", ammindu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammindu.Addr())
	require.Equal(t, 4, ammindu.Len())
	require.Equal(t, uint32(0x3868b5cc), ctorWord(t, ammindu))
}
