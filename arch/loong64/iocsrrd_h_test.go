package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrrdHCtor(t *testing.T) {
	// llvm-mc-verified: iocsrrd.h $t0, $t1.
	require.Equal(
		t,
		uint32(0x064805ac),
		ctorWord(t, NewIocsrrdH(lreg(t, 12), lreg(t, 13))),
	)

	in := NewIocsrrdH(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrrdH)
	require.True(t, ok, "type = %T, want IocsrrdH", in)
}

func TestIocsrrdHDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrrd.h $t0, $t1.
	in := decodeIocsrrdH(0x064805ac, 0x90000000)

	x, ok := in.(IocsrrdH)
	require.True(t, ok, "type = %T, want IocsrrdH", in)
	require.Equal(t, "iocsrrd.h $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064805ac), ctorWord(t, x))
}
