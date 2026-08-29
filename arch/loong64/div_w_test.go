package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestDivWCtor(t *testing.T) {
	// llvm-mc-verified: div.w $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002039ac),
		ctorWord(t, NewDivW(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewDivW(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(DivW)
	require.True(t, ok, "type = %T, want DivW", in)
}

func TestDivWDecodeEncode(t *testing.T) {
	in := decodeDivW(0x002039ac, 0x90000000)

	divw, ok := in.(DivW)
	require.True(t, ok, "type = %T, want DivW", in)
	require.Equal(t, "div.w $t0, $t1, $t2", divw.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), divw.Addr())
	require.Equal(t, uint32(0x002039ac), ctorWord(t, divw))
}
