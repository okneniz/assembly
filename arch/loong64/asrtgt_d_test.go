package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAsrtgtDCtor(t *testing.T) {
	// llvm-mc-verified: asrtgt.d $t1, $t2 (the manual order rj, rk; no rd
	// operand).
	require.Equal(
		t,
		uint32(0x0001b9a0),
		ctorWord(t, NewAsrtgtD(lreg(t, 13), lreg(t, 14))),
	)

	in := NewAsrtgtD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(AsrtgtD)
	require.True(t, ok, "type = %T, want AsrtgtD", in)
}

func TestAsrtgtDDecodeEncode(t *testing.T) {
	in := decodeAsrtgtD(0x0001b9a0, 0x90000000)

	x, ok := in.(AsrtgtD)
	require.True(t, ok, "type = %T, want AsrtgtD", in)
	require.Equal(t, "asrtgt.d $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x0001b9a0), ctorWord(t, x))
}
