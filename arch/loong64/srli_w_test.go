package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSrliWCtor(t *testing.T) {
	// llvm-mc-verified: srli.w $t0, $t1, 3.
	v, err := New().UImm5(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x00448dac),
		ctorWord(t, New().SrliW(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestSrliWDecodeEncode(t *testing.T) {
	x, ok := decodeSrliW(0x00448dac, 0x90000000).(SrliW)
	require.True(t, ok, "type = %T, want SrliW", x)

	require.Equal(t, "srli.w $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x00448dac), ctorWord(t, x))
}
