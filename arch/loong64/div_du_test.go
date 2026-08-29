package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestDivDuCtor(t *testing.T) {
	// llvm-mc-verified: div.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002339ac),
		ctorWord(t, NewDivDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewDivDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(DivDu)
	require.True(t, ok, "type = %T, want DivDu", in)
}

func TestDivDuDecodeEncode(t *testing.T) {
	in := decodeDivDu(0x002339ac, 0x90000000)

	divdu, ok := in.(DivDu)
	require.True(t, ok, "type = %T, want DivDu", in)
	require.Equal(t, "div.du $t0, $t1, $t2", divdu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), divdu.Addr())
	require.Equal(t, uint32(0x002339ac), ctorWord(t, divdu))
}
