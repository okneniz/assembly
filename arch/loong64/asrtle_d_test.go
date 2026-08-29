package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestAsrtleDCtor(t *testing.T) {
	// llvm-mc-verified: asrtle.d $t1, $t2 (the manual order rj, rk; no rd
	// operand).
	require.Equal(
		t,
		uint32(0x000139a0),
		ctorWord(t, NewAsrtleD(lreg(t, 13), lreg(t, 14))),
	)

	in := NewAsrtleD(lreg(t, 1), lreg(t, 2))
	_, ok := in.(AsrtleD)
	require.True(t, ok, "type = %T, want AsrtleD", in)
}

func TestAsrtleDDecodeEncode(t *testing.T) {
	in := decodeAsrtleD(0x000139a0, 0x90000000)

	x, ok := in.(AsrtleD)
	require.True(t, ok, "type = %T, want AsrtleD", in)
	require.Equal(t, "asrtle.d $t1, $t2", x.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), x.Addr())
	require.Equal(t, 4, x.Len())
	require.Equal(t, uint32(0x000139a0), ctorWord(t, x))
}
