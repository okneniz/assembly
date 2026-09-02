package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestPreldxCtor(t *testing.T) {
	// llvm-mc-verified: preldx 5, $t1, $t2 (the manual prints the hint
	// first - llvm-mc rejects the rj-first operand order).
	h, err := New().UImm5(5)
	require.NoError(t, err)

	in := New().Preldx(h, lreg(t, 13), lreg(t, 14))
	require.Equal(t, uint32(0x382c39a5), ctorWord(t, in))

	_, ok := in.(Preldx)
	require.True(t, ok, "type = %T, want Preldx", in)
}

func TestPreldxDecodeEncode(t *testing.T) {
	in := decodePreldx(0x382c39a5, 0x90000000)

	x, ok := in.(Preldx)
	require.True(t, ok, "type = %T, want Preldx", in)
	require.Equal(t, "preldx 5, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.hint.val)
	require.Equal(t, uint32(0x382c39a5), ctorWord(t, x))
}
