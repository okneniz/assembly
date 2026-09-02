package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestIocsrrdDCtor(t *testing.T) {
	// llvm-mc-verified: iocsrrd.d $t0, $t1.
	require.Equal(
		t,
		uint32(0x06480dac),
		ctorWord(t, New().IocsrrdD(lreg(t, 12), lreg(t, 13))),
	)

	in := New().IocsrrdD(lreg(t, 12), lreg(t, 13))
	_, ok := in.(IocsrrdD)
	require.True(t, ok, "type = %T, want IocsrrdD", in)
}

func TestIocsrrdDDecodeEncode(t *testing.T) {
	// llvm-mc-verified: iocsrrd.d $t0, $t1.
	in := decodeIocsrrdD(0x06480dac, 0x90000000)

	x, ok := in.(IocsrrdD)
	require.True(t, ok, "type = %T, want IocsrrdD", in)
	require.Equal(t, "iocsrrd.d $t0, $t1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x06480dac), ctorWord(t, x))
}
