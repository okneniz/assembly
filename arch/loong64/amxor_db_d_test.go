package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmxorDbDCtor(t *testing.T) {
	// llvm-mc-verified: amxor_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386db5cc),
		ctorWord(t, New().AmxorDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmxorDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmxorDbD)
	require.True(t, ok, "type = %T, want AmxorDbD", in)
}

func TestAmxorDbDDecodeEncode(t *testing.T) {
	in := decodeAmxorDbD(0x386db5cc, 0x90000000)

	amxordbd, ok := in.(AmxorDbD)
	require.True(t, ok, "type = %T, want AmxorDbD", in)
	require.Equal(t, "amxor_db.d $t0, $t1, $t2", amxordbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amxordbd.Addr())
	require.Equal(t, 4, amxordbd.Len())
	require.Equal(t, uint32(0x386db5cc), ctorWord(t, amxordbd))
}
