package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrrdBCtor(t *testing.T) {
	// llvm-mc-verified: iocsrrd.b $t0, $t1.
	require.Equal(
		t,
		uint32(0x064801ac),
		ctorWord(t, New().IocsrrdB(lreg(t, 12), lreg(t, 13))),
	)

	in := New().IocsrrdB(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrrdB)
	require.True(t, ok, "type = %T, want IocsrrdB", in)
}

func TestIocsrrdBDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrrd.b $t0, $t1.
	in := decodeIocsrrdB(0x064801ac, 0x90000000)

	x, ok := in.(IocsrrdB)
	require.True(t, ok, "type = %T, want IocsrrdB", in)
	require.Equal(t, "iocsrrd.b $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x064801ac), ctorWord(t, x))
}
