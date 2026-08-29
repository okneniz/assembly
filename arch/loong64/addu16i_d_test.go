package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAddu16iDCtor(t *testing.T) {
	// llvm-mc-verified: addu16i.d $t0, $t1, -1.
	v, err := NewImm16(-1)
	require.NoError(t, err)

	require.Equal(
		t,
		uint32(0x13fffdac),
		ctorWord(t, NewAddu16iD(lreg(t, 12), lreg(t, 13), v)),
	)
}

func TestAddu16iDDecodeEncode(t *testing.T) {
	x, ok := decodeAddu16iD(0x13fffdac, 0x90000000).(Addu16iD)
	require.True(t, ok, "type = %T, want Addu16iD", x)

	require.Equal(t, "addu16i.d $t0, $t1, -1", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, int64(-1), x.imm.val)

	// The negative immediate round-trips through the sign-extended field.
	require.Equal(t, uint32(0x13fffdac), ctorWord(t, x))
}
