package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCsrrdCtor(t *testing.T) {
	csr, err := NewUImm14(5)
	require.NoError(t, err)

	// llvm-mc-verified: csrrd $t0, 5.
	require.Equal(
		t,
		uint32(0x0400140c),
		ctorWord(t, NewCsrrd(lreg(t, 12), csr)),
	)

	in := NewCsrrd(lreg(t, 12), csr)
	_, ok := in.(Csrrd)
	require.True(t, ok, "type = %T, want Csrrd", in)
}

func TestCsrrdDecodeEncode(t *testing.T) {
	// llvm-mc-verified: csrrd $t0, 5.
	in := decodeCsrrd(0x0400140c, 0x90000000)

	x, ok := in.(Csrrd)
	require.True(t, ok, "type = %T, want Csrrd", in)
	require.Equal(t, "csrrd $t0, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0400140c), ctorWord(t, x))
}
