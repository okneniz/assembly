package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmorDbWCtor(t *testing.T) {
	// llvm-mc-verified: amor_db.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386c35cc),
		ctorWord(t, NewAmorDbW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewAmorDbW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmorDbW)
	require.True(t, ok, "type = %T, want AmorDbW", in)
}

func TestAmorDbWDecodeEncode(t *testing.T) {
	in := decodeAmorDbW(0x386c35cc, 0x90000000)

	amordbw, ok := in.(AmorDbW)
	require.True(t, ok, "type = %T, want AmorDbW", in)
	require.Equal(t, "amor_db.w $t0, $t1, $t2", amordbw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amordbw.Addr())
	require.Equal(t, 4, amordbw.Len())
	require.Equal(t, uint32(0x386c35cc), ctorWord(t, amordbw))
}
