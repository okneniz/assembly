package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestPcaddu18iCtor(t *testing.T) {
	// llvm-mc-verified: pcaddu18i $t0, 5.
	v, err := New().Imm20(5)
	require.NoError(t, err)

	in := New().Pcaddu18i(lreg(t, 12), v)
	require.Equal(t, uint32(0x1e0000ac), ctorWord(t, in))

	_, ok := in.(Pcaddu18i)
	require.True(t, ok, "type = %T, want Pcaddu18i", in)
}

func TestPcaddu18iDecodeEncode(t *testing.T) {
	in := decodePcaddu18i(0x1e0000ac, 0x90000000)

	x, ok := in.(Pcaddu18i)
	require.True(t, ok, "type = %T, want Pcaddu18i", in)
	require.Equal(t, "pcaddu18i $t0, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.imm.val)
	require.Equal(t, uint32(0x1e0000ac), ctorWord(t, x))
}
