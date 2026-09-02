package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrwrDCtor(t *testing.T) {
	// llvm-mc-verified: iocsrwr.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x06481dac),
		ctorWord(t, New().IocsrwrD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().IocsrwrD(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrwrD)
	require.True(t, ok, "type = %T, want IocsrwrD", in)
}

func TestIocsrwrDDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrwr.d $t0, $t1.
	in := decodeIocsrwrD(0x06481dac, 0x90000000)

	x, ok := in.(IocsrwrD)
	require.True(t, ok, "type = %T, want IocsrwrD", in)
	require.Equal(t, "iocsrwr.d $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06481dac), ctorWord(t, x))
}
