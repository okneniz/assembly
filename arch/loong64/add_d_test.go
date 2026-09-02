package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAddDCtor(t *testing.T) {
	// llvm-mc-verified: add.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0010b9ac),
		ctorWord(t, New().AddD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().AddD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(AddD)
	require.True(t, ok, "type = %T, want AddD", in)
}

func TestAddDDecodeEncode(t *testing.T) {
	in := decodeOne(0x0010b9ac, 0x90000000)

	x, ok := in.(AddD)
	require.True(t, ok, "type = %T, want AddD", in)
	require.Equal(t, "add.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0010b9ac), ctorWord(t, x))
}
