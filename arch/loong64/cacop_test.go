package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestCacopCtor(t *testing.T) {
	op, err := New().UImm5(5)
	require.NoError(t, err)

	off, err := New().Imm12(8)
	require.NoError(t, err)

	// llvm-mc-verified: cacop 5, $t1, 8 (op, rj, si12).
	require.Equal(
		t,
		uint32(0x060021a5),
		ctorWord(t, New().Cacop(op, lreg(t, 13), off)),
	)

	in := New().Cacop(op, lreg(t, 13), off)
	_, ok := in.(Cacop)
	require.True(t, ok, "type = %T, want Cacop", in)
}

func TestCacopDecodeEncode(t *testing.T) {
	// llvm-mc-verified: cacop 5, $t1, 8.
	in := decodeCacop(0x060021a5, 0x90000000)

	x, ok := in.(Cacop)
	require.True(t, ok, "type = %T, want Cacop", in)
	require.Equal(t, "cacop 5, $t1, 8", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x060021a5), ctorWord(t, x))
}
