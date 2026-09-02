package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAmorDbDCtor(t *testing.T) {
	// llvm-mc-verified: amor_db.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x386cb5cc),
		ctorWord(t, New().AmorDbD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AmorDbD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AmorDbD)
	require.True(t, ok, "type = %T, want AmorDbD", in)
}

func TestAmorDbDDecodeEncode(t *testing.T) {
	in := decodeAmorDbD(0x386cb5cc, 0x90000000)

	amordbd, ok := in.(AmorDbD)
	require.True(t, ok, "type = %T, want AmorDbD", in)
	require.Equal(t, "amor_db.d $t0, $t1, $t2", amordbd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), amordbd.Addr())
	require.Equal(t, 4, amordbd.Len())
	require.Equal(t, uint32(0x386cb5cc), ctorWord(t, amordbd))
}
