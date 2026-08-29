package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestNorCtor(t *testing.T) {
	// data-verified (base | $t0,$t1,$t2): nor $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001439ac),
		ctorWord(t, NewNor(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewNor(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Nor)
	require.True(t, ok, "type = %T, want Nor", in)
}

func TestNorDecodeEncode(t *testing.T) {
	in := decodeOne(0x001439ac, 0x90000000)

	x, ok := in.(Nor)
	require.True(t, ok, "type = %T, want Nor", in)
	require.Equal(t, "nor $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001439ac), ctorWord(t, x))
}
