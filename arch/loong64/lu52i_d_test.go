package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestLu52iDCtor(t *testing.T) {
	// llvm-mc-verified: lu52i.d $t0, $t1, 5.
	v, err := NewImm12(5)
	require.NoError(t, err)

	in := NewLu52iD(lreg(t, 12), lreg(t, 13), v)
	require.Equal(t, uint32(0x030015ac), ctorWord(t, in))

	_, ok := in.(Lu52iD)
	require.True(t, ok, "type = %T, want Lu52iD", in)
}

func TestLu52iDDecodeEncode(t *testing.T) {
	in := decodeLu52iD(0x030015ac, 0x90000000)

	x, ok := in.(Lu52iD)
	require.True(t, ok, "type = %T, want Lu52iD", in)
	require.Equal(t, "lu52i.d $t0, $t1, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.imm.val)
	require.Equal(t, uint32(0x030015ac), ctorWord(t, x))
}
