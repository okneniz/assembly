package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrwrHCtor(t *testing.T) {
	// llvm-mc-verified: iocsrwr.h $t0, $t1.
	require.Equal(
		t,
		uint32(0x064815ac),
		ctorWord(t, New().IocsrwrH(lreg(t, 12), lreg(t, 13))),
	)

	in := New().IocsrwrH(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrwrH)
	require.True(t, ok, "type = %T, want IocsrwrH", in)
}

func TestIocsrwrHDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrwr.h $t0, $t1.
	in := decodeIocsrwrH(0x064815ac, 0x90000000)

	x, ok := in.(IocsrwrH)
	require.True(t, ok, "type = %T, want IocsrwrH", in)
	require.Equal(t, "iocsrwr.h $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064815ac), ctorWord(t, x))
}
