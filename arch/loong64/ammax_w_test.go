package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmmaxWCtor(t *testing.T) {
	// llvm-mc-verified: ammax.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386535cc),
		ctorWord(t, NewAmmaxW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmmaxW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmmaxW)
	require.True(t, ok, "type = %T, want AmmaxW", in)
}

func TestAmmaxWDecodeEncode(t *testing.T) {
	in := decodeAmmaxW(0x386535cc, 0x90000000)

	ammaxw, ok := in.(AmmaxW)
	require.True(t, ok, "type = %T, want AmmaxW", in)
	require.Equal(t, "ammax.w $t0, $t1, $t2", ammaxw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), ammaxw.Addr())
	require.Equal(t, 4, ammaxw.Len())
	require.Equal(t, uint32(0x386535cc), ctorWord(t, ammaxw))
}
