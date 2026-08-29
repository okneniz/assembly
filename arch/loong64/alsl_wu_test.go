package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAlslWuCtor(t *testing.T) {
	// llvm-mc-verified: alsl.wu $t0, $t1, $t2, 3 (ui2 = shift - 1 = 2).
	require.Equal(
		t,
		uint32(0x000739ac),
		ctorWord(t, NewAlslWu(lreg(t, 12), lreg(t, 13), lreg(t, 14), shift3v(t, 3))),
	)

	in := NewAlslWu(lreg(t, 1), lreg(t, 2), lreg(t, 3), shift3v(t, 1))
	_, ok := in.(AlslWu)
	require.True(t, ok, "type = %T, want AlslWu", in)
}

func TestAlslWuDecodeEncode(t *testing.T) {
	in := decodeAlslWu(0x000739ac, 0x90000000)

	x, ok := in.(AlslWu)
	require.True(t, ok, "type = %T, want AlslWu", in)

	// The raw ui2 field is 2; the decoded shift displays field + 1 = 3.
	require.Equal(t, int64(3), x.shift.val)
	require.Equal(t, "alsl.wu $t0, $t1, $t2, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x000739ac), ctorWord(t, x))
}
