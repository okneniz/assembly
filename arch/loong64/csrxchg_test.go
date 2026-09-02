package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCsrxchgCtor(t *testing.T) {
	csr, err := New().UImm14(5)
	require.NoError(t, err)

	// llvm-mc-verified: csrxchg $t0, $t1, 5.
	require.Equal(
		t,
		uint32(0x040015ac),
		ctorWord(t, New().Csrxchg(lreg(t, 12), lreg(t, 13), csr)),
	)

	in := New().Csrxchg(lreg(t, 12), lreg(t, 13), csr)
	_, ok := in.(Csrxchg)
	require.True(t, ok, "type = %T, want Csrxchg", in)
}

func TestCsrxchgDecodeEncode(t *testing.T) {
	// llvm-mc-verified: csrxchg $t0, $t1, 5.
	in := decodeCsrxchg(0x040015ac, 0x90000000)

	x, ok := in.(Csrxchg)
	require.True(t, ok, "type = %T, want Csrxchg", in)
	require.Equal(t, "csrxchg $t0, $t1, 5", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x040015ac), ctorWord(t, x))
}
