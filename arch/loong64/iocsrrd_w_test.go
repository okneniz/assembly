package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrrdWCtor(t *testing.T) {
	// llvm-mc-verified: iocsrrd.w $t0, $t1.
	require.Equal(
		t,
		uint32(0x064809ac),
		ctorWord(t, New().IocsrrdW(lreg(t, 12), lreg(t, 13))),
	)

	in := New().IocsrrdW(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrrdW)
	require.True(t, ok, "type = %T, want IocsrrdW", in)
}

func TestIocsrrdWDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrrd.w $t0, $t1.
	in := decodeIocsrrdW(0x064809ac, 0x90000000)

	x, ok := in.(IocsrrdW)
	require.True(t, ok, "type = %T, want IocsrrdW", in)
	require.Equal(t, "iocsrrd.w $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064809ac), ctorWord(t, x))
}
