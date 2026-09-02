package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestModDuCtor(t *testing.T) {
	// llvm-mc-verified: mod.du $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0023b9ac),
		ctorWord(t, New().ModDu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := New().ModDu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(ModDu)
	require.True(t, ok, "type = %T, want ModDu", in)
}

func TestModDuDecodeEncode(t *testing.T) {
	in := decodeModDu(0x0023b9ac, 0x90000000)

	moddu, ok := in.(ModDu)
	require.True(t, ok, "type = %T, want ModDu", in)
	require.Equal(t, "mod.du $t0, $t1, $t2", moddu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), moddu.Addr())
	require.Equal(t, uint32(0x0023b9ac), ctorWord(t, moddu))
}
