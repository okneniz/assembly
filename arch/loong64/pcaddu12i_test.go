package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestPcaddu12iCtor(t *testing.T) {
	// llvm-mc-verified: pcaddu12i $t0, 5.
	v, err := NewImm20(5)
	require.NoError(t, err)

	in := NewPcaddu12i(lreg(t, 12), v)
	require.Equal(t, uint32(0x1c0000ac), ctorWord(t, in))

	_, ok := in.(Pcaddu12i)
	require.True(t, ok, "type = %T, want Pcaddu12i", in)
}

func TestPcaddu12iDecodeEncode(t *testing.T) {
	in := decodePcaddu12i(0x1c0000ac, 0x90000000)

	x, ok := in.(Pcaddu12i)
	require.True(t, ok, "type = %T, want Pcaddu12i", in)
	require.Equal(t, "pcaddu12i $t0, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, int64(5), x.imm.val)
	require.Equal(t, uint32(0x1c0000ac), ctorWord(t, x))
}
