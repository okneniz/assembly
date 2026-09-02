package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSrlDCtor(t *testing.T) {
	// llvm-mc-verified: srl.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001939ac),
		ctorWord(t, New().SrlD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().SrlD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SrlD)
	require.True(t, ok, "type = %T, want SrlD", in)
}

func TestSrlDDecodeEncode(t *testing.T) {
	in := decodeOne(0x001939ac, 0x90000000)

	x, ok := in.(SrlD)
	require.True(t, ok, "type = %T, want SrlD", in)
	require.Equal(t, "srl.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001939ac), ctorWord(t, x))
}
