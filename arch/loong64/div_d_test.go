package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestDivDCtor(t *testing.T) {
	// llvm-mc-verified: div.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x002239ac),
		ctorWord(t, NewDivD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewDivD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(DivD)
	require.True(t, ok, "type = %T, want DivD", in)
}

func TestDivDDecodeEncode(t *testing.T) {
	in := decodeDivD(0x002239ac, 0x90000000)

	divd, ok := in.(DivD)
	require.True(t, ok, "type = %T, want DivD", in)
	require.Equal(t, "div.d $t0, $t1, $t2", divd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), divd.Addr())
	require.Equal(t, uint32(0x002239ac), ctorWord(t, divd))
}
