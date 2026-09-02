package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCsrwrCtor(t *testing.T) {
	csr, err := New().UImm14(5)
	require.NoError(t, err)

	// llvm-mc-verified: csrwr $t0, 5.
	require.Equal(
		t,
		uint32(0x0400142c),
		ctorWord(t, New().Csrwr(lreg(t, 12), csr)),
	)

	in := New().Csrwr(lreg(t, 12), csr)
	_, ok := in.(Csrwr)
	require.True(t, ok, "type = %T, want Csrwr", in)
}

func TestCsrwrDecodeEncode(t *testing.T) {
	// llvm-mc-verified: csrwr $t0, 5.
	in := decodeCsrwr(0x0400142c, 0x90000000)

	x, ok := in.(Csrwr)
	require.True(t, ok, "type = %T, want Csrwr", in)
	require.Equal(t, "csrwr $t0, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0400142c), ctorWord(t, x))
}
