package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestInvtlbCtor(t *testing.T) {
	op, err := NewUImm5(3)
	require.NoError(t, err)

	// llvm-mc-verified: invtlb 3, $t1, $t2 (op, rj, rk).
	require.Equal(
		t,
		uint32(0x0649b9a3),
		ctorWord(t, NewInvtlb(op, lreg(t, 13), lreg(t, 14))),
	)

	in := NewInvtlb(op, lreg(t, 13), lreg(t, 14))
	_, ok := in.(Invtlb)
	require.True(t, ok, "type = %T, want Invtlb", in)
}

func TestInvtlbDecodeEncode(t *testing.T) {
	// llvm-mc-verified: invtlb 3, $t1, $t2.
	in := decodeInvtlb(0x0649b9a3, 0x90000000)

	x, ok := in.(Invtlb)
	require.True(t, ok, "type = %T, want Invtlb", in)
	require.Equal(t, "invtlb 3, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0649b9a3), ctorWord(t, x))
}
