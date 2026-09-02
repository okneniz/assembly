package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMaskeqzCtor(t *testing.T) {
	// data-verified (base | $t0,$t1,$t2): maskeqz $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001339ac),
		ctorWord(t, New().Maskeqz(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().Maskeqz(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Maskeqz)
	require.True(t, ok, "type = %T, want Maskeqz", in)
}

func TestMaskeqzDecodeEncode(t *testing.T) {
	in := decodeOne(0x001339ac, 0x90000000)

	x, ok := in.(Maskeqz)
	require.True(t, ok, "type = %T, want Maskeqz", in)
	require.Equal(t, "maskeqz $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001339ac), ctorWord(t, x))
}
