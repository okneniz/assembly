package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSubWCtor(t *testing.T) {
	// data-verified (base | $t0,$t1,$t2): sub.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001139ac),
		ctorWord(t, New().SubW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().SubW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SubW)
	require.True(t, ok, "type = %T, want SubW", in)
}

func TestSubWDecodeEncode(t *testing.T) {
	in := decodeOne(0x001139ac)

	x, ok := in.(SubW)
	require.True(t, ok, "type = %T, want SubW", in)
	require.Equal(t, "sub.w $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint32(0x001139ac), ctorWord(t, x))
}
