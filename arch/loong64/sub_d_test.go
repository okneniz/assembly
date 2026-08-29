package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSubDCtor(t *testing.T) {
	// data-verified (base | $t0,$t1,$t2): sub.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0011b9ac),
		ctorWord(t, NewSubD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewSubD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(SubD)
	require.True(t, ok, "type = %T, want SubD", in)
}

func TestSubDDecodeEncode(t *testing.T) {
	in := decodeOne(0x0011b9ac, 0x90000000)

	x, ok := in.(SubD)
	require.True(t, ok, "type = %T, want SubD", in)
	require.Equal(t, "sub.d $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x0011b9ac), ctorWord(t, x))
}
