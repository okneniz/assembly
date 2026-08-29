package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSllWCtor(t *testing.T) {
	// llvm-mc-verified: sll.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001739ac),
		ctorWord(t, NewSllW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewSllW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SllW)
	require.True(t, ok, "type = %T, want SllW", in)
}

func TestSllWDecodeEncode(t *testing.T) {
	in := decodeOne(0x001739ac, 0x90000000)

	x, ok := in.(SllW)
	require.True(t, ok, "type = %T, want SllW", in)
	require.Equal(t, "sll.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001739ac), ctorWord(t, x))
}
