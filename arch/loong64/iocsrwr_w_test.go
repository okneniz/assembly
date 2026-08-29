package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrwrWCtor(t *testing.T) {
	// llvm-mc-verified: iocsrwr.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x064819ac),
		ctorWord(t, NewIocsrwrW(lreg(t, 12), lreg(t, 13))),
	)

	in := NewIocsrwrW(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrwrW)
	require.True(t, ok, "type = %T, want IocsrwrW", in)
}

func TestIocsrwrWDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrwr.w $t0, $t1.
	in := decodeIocsrwrW(0x064819ac, 0x90000000)

	x, ok := in.(IocsrwrW)
	require.True(t, ok, "type = %T, want IocsrwrW", in)
	require.Equal(t, "iocsrwr.w $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064819ac), ctorWord(t, x))
}
