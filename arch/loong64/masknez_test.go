package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestMasknezCtor(t *testing.T) {
	// data-verified (base | $t0,$t1,$t2): masknez $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0013b9ac),
		ctorWord(t, NewMasknez(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewMasknez(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Masknez)
	require.True(t, ok, "type = %T, want Masknez", in)
}

func TestMasknezDecodeEncode(t *testing.T) {
	in := decodeOne(0x0013b9ac, 0x90000000)

	x, ok := in.(Masknez)
	require.True(t, ok, "type = %T, want Masknez", in)
	require.Equal(t, "masknez $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0013b9ac), ctorWord(t, x))
}
