package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSrliDCtor(t *testing.T) {
	// llvm-mc-verified: srli.d $t0, $t1, 3.
	v, err := NewUImm6(3)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x00450dac),
		ctorWord(t, NewSrliD(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestSrliDDecodeEncode(t *testing.T) {
	x, ok := decodeSrliD(0x00450dac, 0x90000000).(SrliD)
	require.True(t, ok, "type = %T, want SrliD", x)

	require.Equal(t, "srli.d $t0, $t1, 3", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(3), x.imm.val)
	require.Equal(t, uint32(0x00450dac), ctorWord(t, x))
}
