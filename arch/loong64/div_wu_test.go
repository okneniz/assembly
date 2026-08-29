package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestDivWuCtor(t *testing.T) {
	// llvm-mc-verified: div.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002139ac),
		ctorWord(t, NewDivWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewDivWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(DivWu)
	require.True(t, ok, "type = %T, want DivWu", in)
}

func TestDivWuDecodeEncode(t *testing.T) {
	in := decodeDivWu(0x002139ac, 0x90000000)

	divwu, ok := in.(DivWu)
	require.True(t, ok, "type = %T, want DivWu", in)
	require.Equal(t, "div.wu $t0, $t1, $t2", divwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), divwu.Addr())
	require.Equal(t, uint32(0x002139ac), ctorWord(t, divwu))
}
