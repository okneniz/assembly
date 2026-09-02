package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSrlWCtor(t *testing.T) {
	// llvm-mc-verified: srl.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0017b9ac),
		ctorWord(t, New().SrlW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().SrlW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SrlW)
	require.True(t, ok, "type = %T, want SrlW", in)
}

func TestSrlWDecodeEncode(t *testing.T) {
	in := decodeOne(0x0017b9ac, 0x90000000)

	x, ok := in.(SrlW)
	require.True(t, ok, "type = %T, want SrlW", in)
	require.Equal(t, "srl.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0017b9ac), ctorWord(t, x))
}
