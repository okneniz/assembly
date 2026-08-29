package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestModDCtor(t *testing.T) {
	// llvm-mc-verified: mod.d $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0022b9ac),
		ctorWord(t, NewModD(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewModD(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(ModD)
	require.True(t, ok, "type = %T, want ModD", in)
}

func TestModDDecodeEncode(t *testing.T) {
	in := decodeModD(0x0022b9ac, 0x90000000)

	modd, ok := in.(ModD)
	require.True(t, ok, "type = %T, want ModD", in)
	require.Equal(t, "mod.d $t0, $t1, $t2", modd.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), modd.Addr())
	require.Equal(t, uint32(0x0022b9ac), ctorWord(t, modd))
}
