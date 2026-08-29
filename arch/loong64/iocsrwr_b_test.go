package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrwrBCtor(t *testing.T) {
	// llvm-mc-verified: iocsrwr.b $t0, $t1.
	require.Equal(
		t,
		uint32(0x064811ac),
		ctorWord(t, NewIocsrwrB(lreg(t, 12), lreg(t, 13))),
	)

	in := NewIocsrwrB(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrwrB)
	require.True(t, ok, "type = %T, want IocsrwrB", in)
}

func TestIocsrwrBDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrwr.b $t0, $t1.
	in := decodeIocsrwrB(0x064811ac, 0x90000000)

	x, ok := in.(IocsrwrB)
	require.True(t, ok, "type = %T, want IocsrwrB", in)
	require.Equal(t, "iocsrwr.b $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064811ac), ctorWord(t, x))
}
