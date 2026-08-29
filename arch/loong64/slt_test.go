package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestSltCtor(t *testing.T) {
	// data-verified (base | $t0,$t1,$t2): slt $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x001239ac),
		ctorWord(t, NewSlt(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewSlt(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(Slt)
	require.True(t, ok, "type = %T, want Slt", in)
}

func TestSltDecodeEncode(t *testing.T) {
	in := decodeOne(0x001239ac, 0x90000000)

	x, ok := in.(Slt)
	require.True(t, ok, "type = %T, want Slt", in)
	require.Equal(t, "slt $t0, $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, uint32(0x001239ac), ctorWord(t, x))
}
