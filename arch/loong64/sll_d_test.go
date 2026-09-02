package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSllDCtor(t *testing.T) {
	// llvm-mc-verified: sll.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0018b9ac),
		ctorWord(t, New().SllD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().SllD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SllD)
	require.True(t, ok, "type = %T, want SllD", in)
}

func TestSllDDecodeEncode(t *testing.T) {
	in := decodeOne(0x0018b9ac, 0x90000000)

	x, ok := in.(SllD)
	require.True(t, ok, "type = %T, want SllD", in)
	require.Equal(t, "sll.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0018b9ac), ctorWord(t, x))
}
