package loong64

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okneniz/assembly/disasm"
)

func TestModWuCtor(t *testing.T) {
	// llvm-mc-verified: mod.wu $t0, $t1, $t2.
	require.Equal(
		t,
		uint32(0x0021b9ac),
		ctorWord(t, NewModWu(lreg(t, 12), lreg(t, 13), lreg(t, 14))),
	)

	in := NewModWu(lreg(t, 1), lreg(t, 2), lreg(t, 3))
	_, ok := in.(ModWu)
	require.True(t, ok, "type = %T, want ModWu", in)
}

func TestModWuDecodeEncode(t *testing.T) {
	in := decodeModWu(0x0021b9ac, 0x90000000)

	modwu, ok := in.(ModWu)
	require.True(t, ok, "type = %T, want ModWu", in)
	require.Equal(t, "mod.wu $t0, $t1, $t2", modwu.ObjDump(disasm.DefaultViewCtx()))
	require.Equal(t, uint64(0x90000000), modwu.Addr())
	require.Equal(t, uint32(0x0021b9ac), ctorWord(t, modwu))
}
